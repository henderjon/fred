//go:build !memory

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func tmp() *os.File {
	f, err := os.CreateTemp("", `fred_tmp_*`)
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func (b bufferLine) String() string {
	return fmt.Sprint(b.pos)
}

type bufferLine struct {
	pos  int
	len  int
	mark rune
}

type scratchBuf struct {
	curline  int
	lastline int
	lines    []bufferLine
	filename string
	ext      io.ReadWriteSeeker
	pos      int
	dirty    bool
}

func newBuffer() buffer {
	f := tmp()
	return &scratchBuf{
		curline:  0,
		lastline: 0,
		lines:    make([]bufferLine, 1),
		filename: "",
		ext:      f,
		dirty:    false,
	}
}

func (b *scratchBuf) getFilename() string {
	return b.filename
}

func (b *scratchBuf) setFilename(fname string) {
	b.filename = fname
}

func (b *scratchBuf) isDirty() bool {
	return b.dirty
}

func (b *scratchBuf) setDirty(d bool) {
	b.dirty = d
}

// Write fulfills io.Writer
func (b *scratchBuf) Write(by []byte) (int, error) {
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
func (b *scratchBuf) Read(p []byte) (int, error) {
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

func (b *scratchBuf) getNumLines() int {
	return b.getLastline()
}

func (b *scratchBuf) setCurline(i int) {
	b.curline = i
}

func (b *scratchBuf) getCurline() int {
	return b.curline
}

func (b *scratchBuf) setLastline(i int) {
	b.lastline = i
}

func (b *scratchBuf) getLastline() int {
	return b.lastline
}

// insertAfter gets input from the user and puts it at the given position
func (b *scratchBuf) insertAfter(inout termio, idx int) error {
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
func (b *scratchBuf) putLine(line string) error {
	b.lastline++

	bts := strings.TrimRightFunc(line, func(r rune) bool {
		return unicode.IsSpace(r)
	})

	num, err := b.ext.Write([]byte(bts))
	if err != nil {
		return err
	}

	newLine := bufferLine{
		pos:  b.pos,
		len:  num,
		mark: null,
	}

	b.pos += num // track tha last bytes written because we'll start there next time
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
func (b *scratchBuf) replaceLine(line string, idx int) error {
	if idx < 1 || idx > b.getLastline() {
		return fmt.Errorf("cannot replace text; invalid address; %d", idx)
	}

	bts := strings.TrimRightFunc(line, func(r rune) bool {
		return unicode.IsSpace(r)
	})

	num, err := b.ext.Write([]byte(bts))
	if err != nil {
		return err
	}

	newLine := bufferLine{
		pos:  b.pos,
		len:  num,
		mark: null,
	}

	b.lines[idx] = newLine
	b.setDirty(true)
	return nil
}

// bulkMove takes the given lines and puts them at dest
func (b *scratchBuf) bulkMove(from, to, dest int) {
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
func (b *scratchBuf) putMark(idx int, r rune) {
	if idx != 0 { // do not allow the zero line to be marked
		b.lines[idx].mark = r
	}
}

// getMark gets the mark of the line at the given index
func (b *scratchBuf) getMark(idx int) rune {
	return b.lines[idx].mark
}

func (b *scratchBuf) hasMark(idx int, r rune) bool {
	return b.getMark(idx) == r
}

// reverse rearranges the given lines in reverse
func (b *scratchBuf) reverse(from, to int) {
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
func (b *scratchBuf) nextLine(n int) int {
	if n >= b.lastline {
		return 0
	}
	return n + 1
}

// prevLine returns the previous index in the buffer, looping to lastline after 0
func (b *scratchBuf) prevLine(n int) int {
	if n <= 0 {
		return b.lastline
	}
	return n - 1
}

// returns the text of the line at the given index
func (b *scratchBuf) getLine(idx int) string {
	b.ext.Seek(int64(b.lines[idx].pos), 0)

	bts := make([]byte, b.lines[idx].len)
	_, err := b.ext.Read(bts)
	if err != nil {
		log.Fatal(err)
	}

	return string(bts)
}

// defLines normalizes two addresses, both optional. It takes what is provided and returns sensible defaults with an eye to how the relate to each other. It also changes '.' and '$' to current and end addresses respectively
func (b *scratchBuf) defLines(start, end, incr string, l1, l2 int) (int, int, error) {
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
			return 0, 0, fmt.Errorf("invalid buffer start address: %s; %s", start, err.Error())
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
			return 0, 0, fmt.Errorf("invalid buffer end address: %s; %s", end, err.Error())
		}
	}

	i1, i2 = b.defIncr(incr, i1, i2)

	if i1 > i2 || i1 <= 0 {
		return 0, 0, fmt.Errorf("invalid buffer range; %d, %d", i1, i2)
	}

	return i1, i2, nil
}

func (b *scratchBuf) defIncr(incr string, start, rel int) (int, int) {
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
func (b *scratchBuf) scanForward(start, num int) func() (int, bool) {
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
func (b *scratchBuf) scanReverse(start, num int) func() (int, bool) {
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

func (b *scratchBuf) destructor() {
	if f, ok := b.ext.(*os.File); ok {
		// log.Println(f.Name())
		os.Remove(f.Name())
	}
}