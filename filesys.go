package main

import (
	"io"
	"os"
	"path/filepath"
)

type fileSystem interface {
	ReadFile(fname string) ([]byte, error)
	WriteFile(data []byte, fname string) (int, error)
	FileReader(fname string) (io.ReadCloser, error)
	FileWriter(fname string) (io.WriteCloser, error)
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (o osFS) ReadFile(fname string) ([]byte, error) {
	absPath, err := o.absPath(fname)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(absPath)
}

func (o osFS) WriteFile(data []byte, fname string) (int, error) {
	absPath, err := o.absPath(fname)
	if err != nil {
		return 0, err
	}

	bts := len(data)
	return bts, os.WriteFile(absPath, data, 0644)
}

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
