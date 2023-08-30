package main

import (
	"fmt"
	"io"
)

const (
	null = '\x00'
	mark = '\x01' // we allow single chars to mark lines, this value is used internally when we go global actions (g,G,v,V) so we don't get confused
)

// buffer is the interface for our content. It is large because of the accessor methods (and the fact that it provides the interface for bufferLines as well) which are not idiomatic, but good practice.
type buffer interface {
	io.ReadWriter
	fmt.Stringer

	defaultLines(start, end, incr string, l1, l2 int) (int, int, error)
	applyIncrement(incr string, start, rel int) (int, int)

	setCurline(int)
	getCurline() int

	setLastline(int)
	getLastline() int

	getFilename() string
	setFilename(string)

	isDirty() bool
	setDirty(bool)
	getRev() int

	putLine(line string, idx int) error
	getLine(int) string

	replaceLine(string, int) error
	duplicateLine(int) error

	bulkMove(from, to, dest int)
	reverse(from, to int)

	putMark(int, rune)
	getMark(int) rune

	nextLine(int) int
	prevLine(int) int

	// these two funcs could be combined but clear is better than clever
	scanForward(int, int) func() (int, bool)
	scanReverse(int, int) func() (int, bool)

	destructor()   // clean up tmp files
	clone() buffer // used for undo
}
