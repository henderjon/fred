package main

import (
	"io"
)

const (
	null = '\x00'
	mark = '\x01' // we allow single chars to mark lines, this value is used internally when we go global actions (g,G,v,V) so we don't get confused
)

// buffer is the interface for our content. It is large because of the accessor methods which are not idiomatic, but good practice.
type buffer interface {
	io.ReadWriter

	defLines(start, end, incr string, l1, l2 int) (int, int, error)
	defIncr(incr string, start, rel int) (int, int)
	getNumLines() int

	setCurline(int)
	getCurline() int

	setLastline(int)
	getLastline() int

	getFilename() string
	setFilename(string)

	isDirty() bool
	setDirty(bool)

	insertAfter(inout termio, idx int) error

	putLine(string) error
	getLine(int) string
	replaceLine(string, int) error

	bulkMove(from, to, dest int)
	reverse(from, to int)

	putMark(int, rune)
	getMark(int) rune
	hasMark(int, rune) bool

	nextLine(int) int
	prevLine(int) int

	// these two funcs could be combined but clear is better than clever
	scanForward(int, int) func() (int, bool)
	scanReverse(int, int) func() (int, bool)

	destructor()
}
