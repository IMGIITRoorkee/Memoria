package memoria

import "io"

// compression interface is a writer and reader which compresses all the data it writes
// and decompresses all the data it reads. the write and reader take a source to do so
type Compression interface {
	Writer(dst io.Writer) (io.WriteCloser, error)
	Reader(src io.Reader) (io.Reader, error)
}

//TODO: Implement compression here
