package main

import (
	"io"
)

func (b bufferLine) String() string {
	return b.txt
}

type bufferLine struct {
	txt  string
	mark bool
}

// buffer is the interface for our content. It is large because of the accessor methods which are not idiomatic, but good practice.
type buffer interface {
	io.ReadWriter

	defLines(start, end string, l1, l2 int) (int, int, error)
	getNumLines() int

	setCurline(i int)
	getCurline() int

	setLastline(i int)
	getLastline() int

	getFilename() string
	setFilename(fname string)

	setPreviousSearch(search)
	getPreviousSearch() search

	insertAfter(input interactor, idx int) error

	putLine(line string) error
	getLine(idx int) string
	replaceLine(line string, idx int) error

	bulkMove(from, to, dest int)
	reverse(from, to int)

	putMark(idx int, m bool)
	getMark(idx int) bool

	nextLine(n int) int
	prevLine(n int) int

	scanForward(int, int) func() (int, bool)
	scanReverse(int, int) func() (int, bool)
}

// TODO: track a drity buffer for save dialog on quit
