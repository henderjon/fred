package main

import "io"

func (b bufferLine) String() string {
	return string(b.txt)
}

type bufferLine struct {
	txt  string
	mark bool
}

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
	setPreviousSearch(pattern string)
	getPreviousSearch() string
	insertAfter(idx int) error
	putText(line string) error
	getLine(idx int) string
	replaceText(line string, idx int) error
	bulkMove(from, to, dest int)
	putMark(idx int, m bool)
	getMark(idx int) bool
	reverse(from, to int)
	nextLine(n int) int
	prevLine(n int) int
}

// TODO: track a drity buffer for save dialog on quit
