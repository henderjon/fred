package main

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type memoryBuf struct {
	curline  int
	lastline int
	lines    []bufferLine
	filename string
}

func newMemoryBuf(fname string) buffer {
	return &memoryBuf{
		curline:  0,
		lastline: 0,
		lines:    make([]bufferLine, 1),
		filename: fname,
	}
}

func (b *memoryBuf) getFilename() string {
	return b.filename
}

func (b *memoryBuf) setFilename(fname string) {
	b.filename = fname
}

func (b *memoryBuf) Write(by []byte) (int, error) {
	var (
		err     error
		line    string
		byCount int
		buf     = bytes.NewBuffer(by)
	)

	for {
		line, err = buf.ReadString('\n')
		if err != nil {
			break
		}

		b.putText(strings.TrimRight(line, "\n\r"))
		byCount += len(line)
	}

	if err == io.EOF {
		err = nil
	}

	return byCount, err
}

func (b *memoryBuf) Read(p []byte) (int, error) {
	var (
		err     error
		byCount int
	)

	for i := 1; i < b.getNumLines(); i++ {
		if byCount >= len(p) {
			return byCount, err
		}

		line := b.getLine(i)
		for _, b := range []byte(line) {
			if byCount >= len(p) {
				return byCount, err
			}
			p = append(p, b)
			byCount++
		}
		// if we read out lines, should we add newlines?
		p = append(p, '\n')
		byCount++
	}
	return byCount, err
}

// func (b *memoryBuf) clear() {
// 	b.curline = 0
// 	b.lastline = 0
// 	b.lines = make([]bufferLine, 1)
// }

func (b *memoryBuf) getNumLines() int {
	return b.getLastline()
	// return len(b.lines) - 1 // take one back for the zero index
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

func (b *memoryBuf) insertAfter(idx int) error {
	var err error

	b.curline = idx
	for stdin.Scan() { // NOTE: use io.Writer
		if err := stdin.Err(); err != nil {
			return err
		}

		line := stdin.Text()

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

// putText adds a new lines to the end of the buffer then moves them into place
func (b *memoryBuf) putText(line string) error {
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

func (b *memoryBuf) getText(idx int) string {
	return b.lines[idx].txt
}

func (b *memoryBuf) replaceText(line string, idx int) error {
	if idx < 1 || idx > b.getLastline() {
		return fmt.Errorf("cannot replace text; invalid address; %d", idx)
	}

	b.lines[idx].txt = line
	return nil
}

func (b *memoryBuf) bulkMove(from, to, dest int) {
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

func (b *memoryBuf) putMark(idx int, m bool) {
	b.lines[idx].mark = m
}

func (b *memoryBuf) getMark(idx int) bool {
	return b.lines[idx].mark
}

func (b *memoryBuf) reverse(from, to int) {
	var tmp bufferLine
	for from < to {
		tmp = b.lines[from]
		b.lines[from] = b.lines[to]
		b.lines[to] = tmp
		from++
		to--
	}
}

func (b *memoryBuf) nextLine(n int) int {
	if n >= b.lastline {
		return 0
	}
	return n + 1
}

func (b *memoryBuf) prevLine(n int) int {
	if n <= 0 {
		return b.lastline
	}
	return n - 1
}

func (b *memoryBuf) getLine(idx int) string {
	return b.lines[idx].String()
}

// defLines normalizes two addresses, both optional. It takes what is provided and returns sensible defaults with an eye to how the relate to each other. It also changes '.' and '$' to current and end addresses respectively
func (b *memoryBuf) defLines(start, end string, l1, l2 int) (int, int, error) {
	var (
		err    error
		i1, i2 int
	)

	if len(start) == 0 && len(end) == 0 {
		return l1, l2, nil
	}

	if start == "." || start == "" { // if no address was given, use the current line
		i1 = b.getCurline()
	} else {
		i1, err = strconv.Atoi(start)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid address (start): %s; %s", start, err.Error())
		}
	}

	switch true {
	case end == "$":
		i2 = b.getLastline()
	case end == "":
		i2 = i1
	default:
		i2, err = strconv.Atoi(end)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid address (end): %s; %s", end, err.Error())
		}
	}

	if i1 > i2 || i1 <= 0 {
		return 0, 0, fmt.Errorf("defaultLines; invalid range; %d, %d", i1, i2)
	}

	return i1, i2, nil
}
