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
	// normalizeAddr(addr string, isSecond bool) int
	defaultLines(start, end string) (int, int, error)
	getNumLines() int
	setCurline(i int)
	insertAfter(idx int, global bool) error
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

// line 1 : integer;
// line2 : integer;
// nlines : integer; number of lines specified
// curln : integer; value of '.'
// lastln : integer; value of '$'
