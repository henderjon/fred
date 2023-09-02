package main

import (
	"io"
	"os"
	"path/filepath"
)

type fileSystem interface {
	Open(name string) (file, error)
}

// defines the file operations used within
type file io.ReadWriteCloser

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Open(name string) (file, error) {
	var err error

	if len(name) == 0 {
		return nil, errEmptyFilename
	}

	absPath, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}

	return os.OpenFile(absPath, os.O_RDWR|os.O_CREATE, 0644)
}
