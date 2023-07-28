package main

type bufferLines []bufferLine

func (b bufferLine) String() string {
	return string(b.txt)
}

type bufferLine struct {
	txt  []byte
	mark bool
}

type buffer interface {
	defaultLines(start, end string) (int, int, error)
	getNumLines() int
	setCurline(i int)
	setLastline(i int)
	getLastline() int

	insertAfter(idx int) error
	putText(line []byte) error
	getText(idx int) []byte
	replaceText(line []byte, idx int) error
	bulkMove(from, to, dest int)
	putMark(idx int, m bool)
	getMark(idx int) bool
	reverse(from, to int)
	nextLine(n int) int
	prevLine(n int) int
	getLine(idx int) string
}

// TODO: track a drity buffer for save dialog on quit
