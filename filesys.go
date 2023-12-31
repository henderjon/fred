package main

import (
	"io"
)

// FileSystem wraps a few simple methods for interacting with an underlying filesystem.
// It serves as an abstraction of either os or in memory actions
type FileSystem interface {
	// Abs check creates an Abs path
	Abs(fname string) (string, error)
	// FileReader wraps the Read and Close methods
	FileReader(fname string) (io.ReadCloser, error)
	// FileWriter wraps the Write and Close methods
	FileWriter(fname string) (io.WriteCloser, error)
	// ScratchFile wraps the ReadAt, WriteAt, Name, and Close methods
	ScratchFile() (NamedScratchFile, error)
}

// NamedScratchFile wraps ReaderAt and WriterAt interfaces as well as the Name() and Close() methods
type NamedScratchFile interface {
	io.ReaderAt
	io.WriterAt
	Name() string
	Close() error
}
