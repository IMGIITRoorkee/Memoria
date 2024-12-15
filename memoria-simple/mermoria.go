package memoria

import (
	"bufio"
	"errors"
	"sync"
)

type PathKey struct {
	Path        []string
	FileName    string
	originalKey string
}

type PathTransform func(key string) *PathKey

type Options struct {
	CacheSize     uint64
	Basedir       string
	PathTransform PathTransform
}

type Writer struct {
	bufWriter  *bufio.Writer
	currentKey string
}

type Reader struct {
	currentKey string
	bufReader  *bufio.Reader
}

type Memoria struct {
	o     Options
	cache map[string][]byte
	w     *Writer
	r     *Reader
	mu    sync.RWMutex
}

func NewWriter() *Writer {

}

func New(options Options) (*Memoria, error) {
	if options.Basedir == "" {
		return nil, errors.New("Base Directory is required")
	}

	w := Writer{
		bufWriter: bufio.NewWriter(),
	}

	return nil, nil

}
