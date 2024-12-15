package memoria

import (
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
	defaultTransform = func(s string) *PathKey { return &PathKey{Path: []string{}, FileName: s} }
)

type PathKey struct {
	Path        []string
	FileName    string
	originalKey string
}

// A Path transform function converts "abcdef" to ["ab","cde","f"]
// so the final  location of the data file will be <basedir>/ab/cde/f/abcdef
type PathTransform func(key string) *PathKey

type Options struct {
	CacheSize uint64
	Basedir   string
	// Tempdir       string solve for issues
	pathPerm      os.FileMode
	filePerm      os.FileMode
	PathTransform PathTransform
	bufferSize    int // the reading and writing is bufferd in memria so this feild represents the size of that buffer
	// compression Compression this field represents a compression mechanism for the store
	// index Indexer this field is for the stores that have some sort of ordering

}
type Memoria struct {
	Options
	cache map[string][]byte
	mu    sync.RWMutex
}

// returns an intiialised Memoria strucutre
func New(o Options) *Memoria {

	if o.Basedir == "" {
		o.Basedir = defaultBaseDir
	}

	if o.PathTransform == nil {
		o.PathTransform = defaultTransform
	}

	if o.bufferSize == 0 {
		o.bufferSize = defaultBufferSize
	}

	if o.CacheSize == 0 {
		o.CacheSize = defaultCacheSize
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

func (m *Memoria) writeWithLock(pathKey *PathKey, r io.Reader, sync bool) error {
	if err := m.createDirIfMissing(pathKey); err != nil {
		return fmt.Errorf("Cannot create directory: %s", err)
	}

	f, err := m.createKeyFile(pathKey)

	if err != nil {
		return fmt.Errorf("cannot create key file: %s", err)
	}

	wc := io.WriteCloser(&nopWriteCloser{})

	//TODO: replace wc with compression writer when implementing compression

	// this is the place where data transfers actually happens when
	// we transfer a read buffer to a writer
	if _, err := io.Copy(wc, r); err != nil {
		f.Close()
		os.Remove(f.Name())
		return fmt.Errorf("error while copying")
	}

	return nil
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
	f, err := os.OpenFile(m.pathFor(pathKey), mode, m.filePerm) //creates the file
	if err != nil {
		return nil, fmt.Errorf("open file: %s", err)
	}
	return f, nil

}

func (m *Memoria) pathFor(pathkey *PathKey) string {
	path := filepath.Join(m.Basedir, filepath.Join(pathkey.Path...))
	return filepath.Join(path, pathkey.FileName)
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
func cleanUp(file *os.File) error {
	if err := file.Close(); err != nil {
		return fmt.Errorf("Cannot close file while cleanup:  %s", err)
	}
	if err := os.Remove(file.Name()); err != nil {
		return fmt.Errorf("Cannot remoave file while cleanup: %s", err)
	}
	return nil
}
