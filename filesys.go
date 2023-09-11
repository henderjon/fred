package main

import (
	"io"
)

type FileSystem interface {
	FileReader(fname string) (io.ReadCloser, error)
	FileWriter(fname string) (io.WriteCloser, error)
	ScratchFile() (NamedScratchFile, error)
}

// NamedReaderWriteAt wraps ReaderAt and WriterAt interfaces as well as the Name() method
type NamedScratchFile interface {
	io.ReaderAt
	io.WriterAt
	Name() string
	Close() error
}
