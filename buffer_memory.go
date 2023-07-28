package main

import (
	"fmt"
	"strconv"
)

type memoryBuf struct {
	curline  int
	lastline int
	lines    []bufferLine
}

func newMemoryBuf() buffer {
	return &memoryBuf{
		curline:  0,
		lastline: 0,
		lines:    make([]bufferLine, 1),
	}
}

func (b *memoryBuf) clear() {
	b.curline = 0
	b.lastline = 0
	b.lines = make([]bufferLine, 1)
}

func (b memoryBuf) getNumLines() int {
	return len(b.lines) - 1 // take one back for the zero index
}

func (b *memoryBuf) setCurline(i int) {
	b.curline = i
}

func (b *memoryBuf) getCurline() int {
	return b.curline
}

func (b *memoryBuf) setLastline(i int) {
	b.lastline = i
}

func (b *memoryBuf) getLastline() int {
	return b.lastline
}

func (b memoryBuf) insertAfter(idx int) error {
	var err error

	b.curline = idx
	for stdin.Scan() { // NOTE: use io.Writer
		if err := stdin.Err(); err != nil {
			return err
		}

		line := stdin.Bytes()

		if len(line) == 1 && line[0] == '.' {
			return nil
		}

		err = b.putText(line)
		if err != nil {
			return err
		}
	}
	return err
}

func (b *memoryBuf) putText(line []byte) error {
	b.lastline++
	newLine := bufferLine{
		txt:  line,
		mark: false,
	}

	// some operations (e.g. `c`) use the last line as scratch space while other simply add new lines
	if b.lastline <= len(b.lines)-1 {
		b.lines[b.lastline] = newLine
	} else {
		b.lines = append(b.lines, newLine)
	}

	b.bulkMove(b.lastline, b.lastline, b.curline)
	b.curline++
	return nil
}

func (b memoryBuf) getText(idx int) []byte {
	return b.lines[idx].txt
}

func (b memoryBuf) replaceText(line []byte, idx int) error {
	// if !hasIdx(currentBuffer, idx) {
	// 	return errAddrOutOfRange
	// }

	b.lines[idx].txt = line
	return nil
}

func (b memoryBuf) bulkMove(from, to, dest int) {
	if dest < from-1 {
		b.reverse(dest+1, from-1)
		b.reverse(from, to)
		b.reverse(dest+1, to)
	} else if dest > to {
		b.reverse(from, to)
		b.reverse(to+1, dest)
		b.reverse(from, dest)
	}
}

func (b memoryBuf) putMark(idx int, m bool) {
	b.lines[idx].mark = m
}

func (b memoryBuf) getMark(idx int) bool {
	return b.lines[idx].mark
}

func (b memoryBuf) reverse(from, to int) {
	var tmp bufferLine
	for from < to {
		tmp = b.lines[from]
		b.lines[from] = b.lines[to]
		b.lines[to] = tmp
		from++
		to--
	}
}

func (b memoryBuf) nextLine(n int) int {
	if n >= b.lastline {
		return 0
	}
	return n + 1
}

func (b memoryBuf) prevLine(n int) int {
	if n <= 0 {
		return b.lastline
	}
	return n - 1
}

func (b memoryBuf) getLine(idx int) string {
	return b.lines[idx].String()
}

// defaultLines normalizes two addresses, both optional. It takes what is provided and returns sensible defaults with an eye to how the relate to each other. It also changes '.' and '$' to current and end addresses respectively
func (b memoryBuf) defaultLines(start, end string) (int, int, error) {
	var err error

	var line1 int
	line1, err = guardAddress(start, b.getCurline(), b.getLastline())
	if err != nil {
		return 0, 0, err
	}

	line2 := line1
	line2, err = guardAddress(end, b.getCurline(), b.getLastline())
	if err != nil {
		return 0, 0, err
	}

	if line1 > line2 || line1 <= 0 {
		// we might get a "0" from the command, let them know we don't like that
		return 0, 0, fmt.Errorf("defaultLines; invalid range; %d, %d", line1, line2)
	}

	return line1, line2, nil // page 188
}

// converts a string address into a number with special cases for '.' and '$'
// func (b memoryBuf) defaultLine(addr string) (int, error) {
// 	if addr == "." { // || addr == "" ... we do this check above
// 		return b.curline, nil
// 	}

// 	if addr == "$" {
// 		return b.lastline, nil
// 	}

// 	return strconv.Atoi(addr)
// }

// converts a string address into a number with special cases for '.', '$', and ”. Start/end addresses are guarded against '0' elsewhere (in defaultLines) but allowed in destinations
func guardAddress(addr string, current, last int) (int, error) {
	if addr == "." || addr == "" { // if no address was given, use the current line
		return current, nil
	}

	if addr == "$" {
		return last, nil
	}

	i, err := strconv.Atoi(addr)
	if err != nil {
		return 0, fmt.Errorf("invalid address: %s; %s", addr, err.Error())
	}

	if i < 0 || i > last {
		return 0, fmt.Errorf("invalid address: %s", addr)
	}

	return i, nil
}
