package main

import (
	"bytes"
	"io"
)

// localFS implements fileSystem using the embedded files
type localFS struct {
	seed *localFile
}

func (m *localFS) Abs(fname string) (string, error) {
	return fname, nil
}

func (m *localFS) FileReader(fname string) (io.ReadCloser, error) {
	return m.seed, nil
}

func (m *localFS) FileWriter(fname string) (io.WriteCloser, error) {
	return m.seed, nil
}

type localFile struct {
	bytes.Buffer
}

func (localFile) Close() error {
	return nil
}

func (m *localFS) ScratchFile() (NamedScratchFile, error) {
	return nil, nil
}
