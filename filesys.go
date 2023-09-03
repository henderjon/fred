package main

import (
	"io"
	"os"
	"path/filepath"
)

type fileSystem interface {
	FileReader(fname string) (io.ReadCloser, error)
	FileWriter(fname string) (io.WriteCloser, error)
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) absPath(fname string) (string, error) {
	if len(fname) == 0 {
		return "", errEmptyFilename
	}

	return filepath.Abs(fname)
}

func (o osFS) FileReader(fname string) (io.ReadCloser, error) {
	absPath, err := o.absPath(fname)
	if err != nil {
		return nil, err
	}

	return os.Open(absPath)
}

func (o osFS) FileWriter(fname string) (io.WriteCloser, error) {
	var err error

	absPath, err := o.absPath(fname)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	f.Truncate(0)
	f.Seek(0, 0)

	return f, err
}
