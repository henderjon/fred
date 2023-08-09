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
	search   search
}

func newMemoryBuf() buffer {
	return &memoryBuf{
		curline:  0,
		lastline: 0,
		lines:    make([]bufferLine, 1),
		filename: "",
		search:   search{},
	}
}

func (b *memoryBuf) getFilename() string {
	return b.filename
}

func (b *memoryBuf) setFilename(fname string) {
	b.filename = fname
}

// Write fulfills io.Writer
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

		b.putLine(strings.TrimRight(line, "\n\r"))
		byCount += len(line)
	}

	if err == io.EOF {
		err = nil
	}

	return byCount, err
}

// Read fulfills io.Reader
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

// insertAfter gets input from the user and puts it at the given position
func (b *memoryBuf) insertAfter(input interactor, idx int) error {
	b.curline = idx
	for {
		line, err := input("")
		if err != nil {
			return err
		}

		if len(line) == 1 && line[0] == '.' {
			return nil
		}

		err = b.putLine(line)
		if err != nil {
			return err
		}
	}
}

// putLine adds a new lines to the end of the buffer then moves them into place
func (b *memoryBuf) putLine(line string) error {
	b.lastline++
	newLine := bufferLine{
		txt:  line,
		mark: null,
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

// replaceLine changes the line to the new text at the given index
func (b *memoryBuf) replaceLine(line string, idx int) error {
	if idx < 1 || idx > b.getLastline() {
		return fmt.Errorf("cannot replace text; invalid address; %d", idx)
	}

	b.lines[idx].txt = line
	return nil
}

// bulkMove takes the given lines and puts them at dest
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

// putMark sets the mark of the line at the given index
func (b *memoryBuf) putMark(idx int, r rune) {
	b.lines[idx].mark = r
}

// getMark gets the mark of the line at the given index
func (b *memoryBuf) getMark(idx int) rune {
	return b.lines[idx].mark
}

func (b *memoryBuf) hasMark(idx int, r rune) bool {
	return b.getMark(idx) == r
}

// reverse rearranges the given lines in reverse
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

// nextLine returns the next index in the buffer, looping to 0 after lastline
func (b *memoryBuf) nextLine(n int) int {
	if n >= b.lastline {
		return 0
	}
	return n + 1
}

// prevLine returns the previous index in the buffer, looping to lastline after 0
func (b *memoryBuf) prevLine(n int) int {
	if n <= 0 {
		return b.lastline
	}
	return n - 1
}

// returns the text of the line at the given index
func (b *memoryBuf) getLine(idx int) string {
	return b.lines[idx].txt
}

// store the last successful search pattern
func (b *memoryBuf) setPreviousSearch(search search) {
	b.search = search
}

// get the last successful search pattern
func (b *memoryBuf) getPreviousSearch() search {
	return b.search
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

	switch true {
	case start == "": // if no address was given, use the current line
		i1 = b.getCurline()
	case start == ".": // if no address was given, use the current line
		i1 = b.getCurline()
	case start == "$": // if no address was given, use the current line
		i1 = b.getLastline()
	default:
		i1, err = strconv.Atoi(start)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid buffer address (start): %s; %s", start, err.Error())
		}
	}

	switch true {
	case end == "$":
		i2 = b.getLastline()
	case end == ".":
		i2 = b.getCurline()
	case end == "":
		i2 = i1
	default:
		i2, err = strconv.Atoi(end)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid buffer address (end): %s; %s", end, err.Error())
		}
	}

	if i1 > i2 || i1 <= 0 {
		return 0, 0, fmt.Errorf("invalid buffer range; %d, %d", i1, i2)
	}

	return i1, i2, nil
}

// scanForward returns a func that walks the buffer's indices in a forward loop. As an implementation detail, the number of lines is the number of non-zero lines.
func (b *memoryBuf) scanForward(start, num int) func() (int, bool) {
	stop := false
	i := b.prevLine(start) // remove 1 because nextLine advances one

	num = func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}(num)

	n := 0
	return func() (int, bool) {
		if n <= num { // if we loop back, we'll return '0' so have to stop a '>' and not '>="
			n++
			i = b.nextLine(i)
			return i, !stop
		}

		return i, stop
	}
}

// scanReverse returns a func that walks the buffer's indices in a reverse loop
func (b *memoryBuf) scanReverse(start, num int) func() (int, bool) {
	stop := false
	i := b.nextLine(start) // remove 1 because nextLine advances one

	num = func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}(num)

	n := 0
	return func() (int, bool) {
		if n <= num { // if we loop back, we'll return '0' so have to stop a '>' and not '>="
			n++
			i = b.prevLine(i)
			return i, !stop
		}

		return i, stop
	}
}

// i := start - 1
// return func() (int, bool) {
// 	i = b.nextLine(i)
// 	if (start == end || end == -1) && i != start-1 { // full loop
// 		return i, false
// 	}

// 	if (start != end && end != -1) && i == end {
// 		return i, false
// 	}

// 	return i, true
// }
