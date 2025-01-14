package memoria

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const (
	defaultBufferSize             = 4096 // 4kb default buffer size
	defaultPathPerm   os.FileMode = 0777
	defaultFilePerm   os.FileMode = 0666
	defaultBaseDir                = "memoria"
	defaultCacheSize              = 512 // 512 bytes as default cache size
)

var (
	defaultTransform = func(s string) *PathKey {
		return &PathKey{Path: []string{}, FileName: s}
	}
)
var defaultInverseTransform = func(pathKey *PathKey) string {
	// Rebuild the key by joining the path parts and appending the filename
	return fmt.Sprintf("%s/%s", filepath.Join(pathKey.Path...), pathKey.FileName)
}

type PathKey struct {
	Path        []string
	FileName    string
	originalKey string
}

// A Path transform function converts "abcdef" to ["ab","cde","f"]
// so the final  location of the data file will be <basedir>/ab/cde/f/abcdef
type PathTransform func(key string) *PathKey

// inv transform func takes file path as input and returns the original key
type InversePathTransform func(pathKey *PathKey) string

type Options struct {
	MaxCacheSize uint64
	Basedir      string
	// Tempdir       string solve for issues
	pathPerm             os.FileMode
	filePerm             os.FileMode
	PathTransform        PathTransform
	InversePathTransform InversePathTransform
	cachePolicy          CachePolicy
	bufferSize           int // the reading and writing is bufferd in memria so this feild represents the size of that buffer
	// compression Compression this field represents a compression mechanism for the store
	// index Indexer this field is for the stores that have some sort of ordering

}
type Memoria struct {
	Options
	cache     map[string][]byte
	mu        sync.RWMutex
	cacheSize uint64
}

// returns an intiialised Memoria strucutre
func New(o Options) *Memoria {

	if o.Basedir == "" {
		o.Basedir = defaultBaseDir
	}

	if o.PathTransform == nil {
		o.PathTransform = defaultTransform
	}

	if o.InversePathTransform == nil {
		o.InversePathTransform = defaultInverseTransform
	}

	if o.cachePolicy == nil {
		o.cachePolicy = &defaultCachePolicy{}
	}

	if o.bufferSize == 0 {
		o.bufferSize = defaultBufferSize
	}

	if o.MaxCacheSize == 0 {
		o.MaxCacheSize = defaultCacheSize
	}

	if o.filePerm == 0 {
		o.filePerm = defaultFilePerm
	}

	if o.pathPerm == 0 {
		o.filePerm = defaultPathPerm
	}

	m := &Memoria{
		Options: o,
		cache:   make(map[string][]byte),
	}
	return m
}

func (m *Memoria) transform(key string) (pathkey *PathKey) {
	pathkey = m.PathTransform(key)
	pathkey.originalKey = key
	return pathkey
}

func (m *Memoria) InverseTransform(pathKey *PathKey) string {
	return m.InversePathTransform(pathKey)
}

// Write synchronously the key-value pair to the disk making it immedialtely avaialble for
// reads. If you need stronger sync gaurantess see WriteStream
func (m *Memoria) Write(key string, val []byte) error {
	return m.WriteStream(key, bytes.NewReader(val), false, false)
}

func (m *Memoria) WriteWithAppend(key string, val []byte) error {
	return m.WriteStream(key, bytes.NewReader(val), true, false)
}

// WriteString is a wrapper for Write that takes a string and writes it as bytes
func (m *Memoria) WriteString(key, val string) error {
	return m.Write(key, []byte(val))
}

// writes the data given by the io.reader  performs explicit sync if mentioned otherwise
// depedning on the physical media it sync
func (m *Memoria) WriteStream(key string, r io.Reader, append bool, sync bool) error { //adding the append bool

	if len(key) <= 0 {
		return fmt.Errorf("Empty key")
	}

	pathKey := m.transform(key)

	//TODO: check for bad paths check if any part contains / after being transformed

	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.createDirIfMissing(pathKey); err != nil {
		return fmt.Errorf("Cannot create directory: %s", err)
	}

	if !append {

		f, err := m.createKeyFile(pathKey)

		if err != nil {
			return fmt.Errorf("cannot create key file: %s", err)
		}

		wc := io.WriteCloser(&nopWriteCloser{f})

		//TODO: replace wc with compression writer when implementing compression

		// this is the place where data transfers actually happens when
		// we transfer a read buffer to a writer
		if _, err := io.Copy(wc, r); err != nil {
			return cleanUp(f, fmt.Errorf("Cannot copy from read buffer %s", err))
		}

		if err := wc.Close(); err != nil {
			return cleanUp(f, fmt.Errorf("Cannot close compression error %s", err))
		}

		if sync {
			if err := f.Sync(); err != nil {
				cleanUp(f, fmt.Errorf("Cannot Sync: %s", err))
			}
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("Cannot close file: %s", err)
		}

		//Atomic Writes: uncomment the following code when implemented atomic writes
		// fullPath := m.completePath(pathKey)

		// if f.Name() != fullPath {
		// 	if err := os.Rename(f.Name(), fullPath); err != nil {
		// 		os.Remove(f.Name())
		// 		return fmt.Errorf("Cannot rename files: %s", err)
		// 	}
		// }

		// empty the cache for original key
		m.emptyCacheFor(pathKey.originalKey) // cache is read only
	}

	if append {
		f, err := m.createKeyFileWithAppend(pathKey, append)
		if err != nil {
			return fmt.Errorf("cannot create key file: %s", err)
		}

		// Use the writer directly, handling compression or other transformations here if necessary
		wc := io.WriteCloser(&nopWriteCloser{f})

		// Perform the data copy operation
		if _, err := io.Copy(wc, r); err != nil {
			return cleanUp(f, fmt.Errorf("Cannot copy from read buffer: %s", err))
		}

		// Close the write closer
		if err := wc.Close(); err != nil {
			return cleanUp(f, fmt.Errorf("Cannot close file after writing: %s", err))
		}

		// Sync if required
		if sync {
			if err := f.Sync(); err != nil {
				cleanUp(f, fmt.Errorf("Cannot sync file: %s", err))
			}
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("Cannot close file after sync: %s", err)
		}

		// Empty cache after write if necessary
		m.emptyCacheFor(pathKey.originalKey)

	}

	return nil

}

func (m *Memoria) emptyCacheFor(key string) {
	if val, ok := m.cache[key]; ok {
		m.cacheSize -= uint64(len(val))
		delete(m.cache, key)
	}
}

func (m *Memoria) createDirIfMissing(pathkey *PathKey) error {
	return os.MkdirAll(m.pathFor(pathkey), m.pathPerm)
}

func (m *Memoria) createKeyFile(pathKey *PathKey) (*os.File, error) {
	// if m.Tempdir != "" {
	//TODO: Write implementation here
	// }
	mode := os.O_CREATE | os.O_WRONLY | os.O_TRUNC // defines the mode of operation for creating the file
	// O_WRONLY: Open for writing only
	// O_CREATE: Create the file if it does not exist
	// O_TRUNC: if file exists truncate it to length 0
	f, err := os.OpenFile(m.completePath(pathKey), mode, m.filePerm) //creates the file
	if err != nil {
		return nil, fmt.Errorf("open file: %s", err)
	}
	return f, nil

}

func (m *Memoria) createKeyFileWithAppend(pathKey *PathKey, append bool) (*os.File, error) {
	var mode int
	if append {
		mode = os.O_APPEND | os.O_WRONLY
	} else {
		mode = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	f, err := os.OpenFile(m.completePath(pathKey), mode, m.filePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %s", m.completePath(pathKey), err)
	}

	return f, nil
}

func (m *Memoria) Read(key string) ([]byte, error) {
	rc, err := m.ReadStream(key, false)
	if err != nil {
		return []byte{}, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// ReadString is a wrapper around Read that returns the value as a string
func (m *Memoria) ReadString(key string) (string, error) {
	val, err := m.Read(key)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

// ReadStream takes the key and a bool byPassCache to bypass the cache and laziliy
// delete all the contents of cache for the hit

func (m *Memoria) ReadStream(key string, bypassCache bool) (io.ReadCloser, error) {
	pathKey := m.transform(key)
	m.mu.RLock()
	defer m.mu.RUnlock()

	if val, ok := m.cache[key]; ok {
		if !bypassCache {
			buf := bytes.NewReader(val)
			//COMPRESSION: make this the compression reader in case of compression
			return io.NopCloser(buf), nil
		}
		go func() {
			m.mu.Lock()
			defer m.mu.Unlock()
			m.cacheSize -= uint64(len(val))
			delete(m.cache, key)
		}()
	}

	// read the file from disk in case of cache miss or bypass cache
	fileName := m.completePath(pathKey)

	//TODO: first check whether filename is valid or not to return appropiate error
	// use os.Stat
	f, err := os.Open(fileName)

	if err != nil {
		return nil, fmt.Errorf("Cannot open file %s", err)
	}

	var r io.Reader

	if m.MaxCacheSize > 0 {
		r = newCachingReader(f, m, key)
	} else {
		r = &closingReader{f}
	}

	var rc = io.ReadCloser(io.NopCloser(r))

	return rc, nil
}

// closingReader provides a Reader that automatically closes the
// embedded ReadCloser when it reaches EOF
type closingReader struct {
	rc io.ReadCloser
}

func (cr closingReader) Read(p []byte) (int, error) {
	n, err := cr.rc.Read(p)
	if err == io.EOF {
		if closeErr := cr.rc.Close(); closeErr != nil {
			return n, closeErr // close must succeed for Read to succeed
		}
	}
	return n, err
}

// this denotes a reader which also caches the data as it reads this in case when size
// of the cache is greater than 0
type cachingReader struct {
	f   *os.File
	m   *Memoria
	key string
	buf *bytes.Buffer
}

func newCachingReader(f *os.File, m *Memoria, key string) io.Reader {
	return &cachingReader{
		f:   f,
		m:   m,
		key: key,
		buf: &bytes.Buffer{},
	}
}

// read interface for io.Reader
func (c *cachingReader) Read(p []byte) (int, error) {
	n, err := c.f.Read(p)

	if err == nil {
		return c.buf.Write(p[0:n]) // write must succedd for read to succed
	}

	if err == io.EOF {
		if err := c.m.cacheWithoutLock(c.key, c.buf.Bytes()); err != nil {
			return n, err
		} // cache may fail

		if closeErr := c.f.Close(); closeErr != nil {
			return n, closeErr
		}
	}

	return n, err

}

func (m *Memoria) pathFor(pathkey *PathKey) string {
	return filepath.Join(m.Basedir, filepath.Join(pathkey.Path...))
}
func (m *Memoria) completePath(path *PathKey) string {
	return filepath.Join(m.pathFor(path), path.FileName)
}

// cache the give key-value pain
func (m *Memoria) cacheWithLock(key string, val []byte) error {
	m.emptyCacheFor(key) // remove the cache if it already exists

	valueSize := uint64(len(val))

	if err := m.makeSpace(valueSize); err != nil {
		return fmt.Errorf("%s; cannot cache", err)
	}

	if err := m.cachePolicy.Insert(m, key, val); err != nil {
		return fmt.Errorf("%s; cannot insert", err)
	}
	return nil
}

func (m *Memoria) makeSpace(valueSize uint64) error {
	if valueSize > m.MaxCacheSize {
		return fmt.Errorf("value size (%d bytes) is too large for cache (%d bytes)", valueSize, m.MaxCacheSize)
	}
	// how much space we need
	spaceNeeded := (m.cacheSize + valueSize) - m.MaxCacheSize
	return m.cachePolicy.Eject(m, spaceNeeded)
}

// aquires the store's mutex and calls Lock
func (m *Memoria) cacheWithoutLock(key string, val []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cacheWithLock(key, val)
}

// If you do compression you need both write and close interfaces so the
// file writes are generic with help of writer and closer
// This is a no-op i.e no operation closer
type nopWriteCloser struct {
	io.Writer
}

func (wc *nopWriteCloser) Write(p []byte) (int, error) { return wc.Writer.Write(p) }
func (wc *nopWriteCloser) Close() error                { return nil }

// /// HELPER FUNCTIONS REFACTOR PLEASE!
func cleanUp(file *os.File, onCleanUpError error) error {
	if err := file.Close(); err != nil {
		return fmt.Errorf("Cannot close file while cleanup:  %s", err)
	}
	if err := os.Remove(file.Name()); err != nil {
		return fmt.Errorf("Cannot remoave file while cleanup: %s", err)
	}
	return fmt.Errorf("%s ..Files Cleaned!", onCleanUpError)
}

// Implementing Concurrent Bulk Write Operations using Go Routines

type WriteResult struct {
	Key   string
	Error error
}

func (m *Memoria) BulkWrite(pairs map[string][]byte, numWorkers int) []WriteResult { // I've taken the number of workers as an argument so that we've a control over how many goroutines we wanna use

	var wg sync.WaitGroup
	results := make([]WriteResult, 0, len(pairs)) //To store results of each write op and also I've kept its size equal to no. of pairs

	// Channel for Worker Pool
	workChan := make(chan struct {
		key   string
		value []byte
	}, len(pairs))

	//Creating channel for goroutines
	resultChan := make(chan WriteResult, len(pairs)) //Hence, the Buffer size in channel is equal to the no. of pairs

	// Worker pool
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range workChan {
				err := m.Write(task.key, task.value)
				resultChan <- WriteResult{
					Key:   task.key,
					Error: err,
				}
			}
		}()
	}

	// Feeding tasks to the work channel
	go func() {
		for key, value := range pairs {
			workChan <- struct {
				key   string
				value []byte
			}{key, value}
		}
		close(workChan) // To close the work channel once all tasks are fed
	}()

	go func() {
		wg.Wait()
		close(resultChan) //To close the channel once all the goroutines are completed
	}()

	for result := range resultChan {
		results = append(results, result)
	}

	return results

}

// Implementing the Close() method:

func (m *Memoria) Close() error { // Here, the Close() method only clears the in-memory cache
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clearing the cache within the memory:
	for key := range m.cache {
		delete(m.cache, key) // To delete the key from cache map
	}
	m.cacheSize = 0

	return nil
}

// createdump method
func (m *Memoria) createDump() error {
	// buffer to hold the dump data
	var buf bytes.Buffer

	// writer to serialize the key-value data into the buffer.
	encoder := json.NewEncoder(&buf)

	m.mu.RLock()
	defer m.mu.RUnlock()

	for key, value := range m.cache {
		// Write key-value pairs into the buffer
		if err := encoder.Encode(map[string][]byte{key: value}); err != nil {
			return fmt.Errorf("failed to encode key-value pair %s: %v", key, err)
		}
	}

	dumpFilePath := filepath.Join(m.Basedir, "backup.dump")
	file, err := os.Create(dumpFilePath)
	if err != nil {
		return fmt.Errorf("failed to create dump file: %v", err)
	}
	defer file.Close()

	if _, err := buf.WriteTo(file); err != nil {
		return fmt.Errorf("failed to write dump data to file: %v", err)
	}

	return nil
}

// Backup mrthod restores the store's data from a backup file located in the given directory.
func (m *Memoria) Backup(backupDir string) error {
	// Construct path to the dump file that needs to be restored
	dumpFilePath := filepath.Join(backupDir, "backup.dump")

	// Open the backup dump file
	file, err := os.Open(dumpFilePath)
	if err != nil {
		return fmt.Errorf("failed to open backup dump file: %v", err)
	}
	defer file.Close()

	// Read the contents of the file into a buffer
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(file); err != nil {
		return fmt.Errorf("failed to read from backup dump file: %v", err)
	}

	// Decode the backup data into memory
	decoder := json.NewDecoder(&buf)
	for {
		var data map[string][]byte
		if err := decoder.Decode(&data); err != nil {
			if err == io.EOF {
				break // End of the backup file
			}
			return fmt.Errorf("failed to decode backup data: %v", err)
		}

		// Insert each key-value pair back into the memory store
		for key, value := range data {
			m.mu.Lock()
			m.cache[key] = value
			m.mu.Unlock()
		}
	}

	return nil
}
