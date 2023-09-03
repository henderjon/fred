package main

import (
	"os"
	"path/filepath"
)

type fileSystem interface {
	ReadFile(fname string) ([]byte, error)
	WriteFile(data []byte, fname string) (int, error)
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) ReadFile(fname string) ([]byte, error) {
	var err error

	if len(fname) == 0 {
		return nil, errEmptyFilename
	}

	absPath, err := filepath.Abs(fname)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(absPath)
}

func (osFS) WriteFile(data []byte, fname string) (int, error) {
	bts := len(data)
	return bts, os.WriteFile(fname, data, 0644)
}
