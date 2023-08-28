//go:build memory

package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type bufferLine struct {
	txt  string
	mark rune
}

func (b bufferLine) String() string {
	return fmt.Sprintf("len: %d; mark: %c; %s", len(b.txt), b.mark, b.txt)
}

type memoryBuf struct {
	curline  int
	lastline int
	lines    []bufferLine
	filename string
	dirty    bool
	rev      int
}

func newBuffer() buffer {
	return &memoryBuf{
		curline:  0,
		lastline: 0,
		lines:    make([]bufferLine, 1),
		filename: "",
		dirty:    false,
		rev:      0,
	}
}

func (b *memoryBuf) getFilename() string {
	return b.filename
}

func (b *memoryBuf) setFilename(fname string) {
	b.filename = fname
}

func (b *memoryBuf) isDirty() bool {
	return b.dirty
}

func (b *memoryBuf) setDirty(d bool) {
	if d {
		b.rev++
	}
	b.dirty = d
}

func (b *memoryBuf) getRev() int {
	return b.rev
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
func (b *memoryBuf) insertAfter(inout termio, idx int) error {
	b.curline = idx
	for {
		line, err := inout.input("")
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
		txt: strings.TrimRightFunc(line, func(r rune) bool {
			return unicode.IsSpace(r)
		}),
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
	b.setDirty(true)
	return nil
}

// duplicateLine changes the line to the new text at the given index
func (b *memoryBuf) duplicateLine(idx int) error {
	if idx < 1 || idx > b.getLastline() {
		return fmt.Errorf("cannot duplicate text; invalid address; %d", idx)
	}

	b.lastline++

	// if there are rando [forgotten] lines at the end of the buffer, reuse them
	if b.lastline <= len(b.lines)-1 {
		b.lines[b.lastline] = b.lines[idx]
	} else {
		b.lines = append(b.lines, b.lines[idx])
	}

	// NOTE: an argument can be made to do the bulk move here, one at a time and empty doCopyNPaste ...

	b.setDirty(true)
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
	b.setDirty(true) // sometimes neither if clause is satisfied
}

// putMark sets the mark of the line at the given index
func (b *memoryBuf) putMark(idx int, r rune) {
	if idx != 0 { // do not allow the zero line to be marked
		b.lines[idx].mark = r
	}
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
	b.setDirty(true)
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

// defLines normalizes two addresses, both optional. It takes what is provided and returns sensible defaults with an eye to how the relate to each other. It also changes '.' and '$' to current and end addresses respectively
func (b *memoryBuf) defLines(start, end, incr string, l1, l2 int) (int, int, error) {
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
		i1, err = intval(start)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid buffer start address: %s", err.Error())
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
		i2, err = intval(end)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid buffer end address: %s", err.Error())
		}
	}

	i1, i2 = b.defIncr(incr, i1, i2)

	if i1 > i2 || i1 <= 0 {
		return 0, 0, fmt.Errorf("invalid buffer range; %d, %d", i1, i2)
	}

	return i1, i2, nil
}

func (b *memoryBuf) defIncr(incr string, start, rel int) (int, int) {
	end := rel
	if incr == ">" {
		end = start + rel
	}

	if incr == "<" {
		end = start
		start -= rel
	}

	if start < 1 {
		start = 1
	}

	if end > b.getLastline() {
		end = b.getLastline()
	}

	return start, end

}

// scanForward returns a func that walks the buffer's indices in a forward loop. As an implementation detail, the number of lines is the number of non-zero lines.
func (b *memoryBuf) scanForward(start, num int) func() (int, bool) {
	stop := false
	i := b.prevLine(start) // remove 1 because nextLine advances one

	if num < 0 {
		num = b.getNumLines()
	}

	n := 0 // this is essentially a do{}while() loop so '0' will execute once
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

	if num < 0 {
		num = b.getNumLines()
	}

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

func (b *memoryBuf) destructor() {}

func (b *memoryBuf) clone() buffer {
	cll := make([]bufferLine, len(b.lines))
	copy(cll, b.lines)
	t := &memoryBuf{
		curline:  b.curline,
		lastline: b.lastline,
		lines:    cll,
		filename: b.filename,
		dirty:    b.dirty,
		rev:      b.rev,
	}
	return t
}

func (b *memoryBuf) String() string {
	var rtn strings.Builder
	fmt.Fprint(&rtn, "buffer (memory):\r\n")
	fmt.Fprintf(&rtn, "  curline: %d\r\n", b.curline)
	fmt.Fprintf(&rtn, "  lastline: %d\r\n", b.lastline)
	fmt.Fprintf(&rtn, "  filename: %s\r\n", b.filename)
	fmt.Fprintf(&rtn, "  dirty: %t\r\n", b.dirty)
	fmt.Fprintf(&rtn, "  rev: %d\r\n", b.rev)
	for k, v := range b.lines {
		if k == 0 {
			continue
		}
		fmt.Fprintf(&rtn, "  line[%d]: %s\r\n", k, v.String())
	}
	return rtn.String()
}
