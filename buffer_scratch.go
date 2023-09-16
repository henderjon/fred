//go:build !memory

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type bufferLine struct {
	pos  int
	len  int
	mark rune
}

func (b bufferLine) String() string {
	return fmt.Sprintf("pos: %d; len: %d; mark: %c %d;", b.pos, b.len, b.mark, b.mark)
}

type scratchBuf struct {
	curline  int
	lastline int
	lines    []bufferLine
	filename string
	ext      NamedScratchFile
	pos      int
	dirty    bool
	rev      int // revision, each time we alter a buffer we incr
}

func newBuffer(fs FileSystem) buffer {
	f, e := fs.ScratchFile()
	if e != nil {
		stderr.Fatal(e)
	}
	return &scratchBuf{
		curline:  0,
		lastline: 0,
		lines:    make([]bufferLine, 1),
		filename: "",
		ext:      f,
		dirty:    false,
		rev:      0,
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
	if d {
		b.rev++
	}
	b.dirty = d
}

func (b *scratchBuf) getRev() int {
	return b.rev
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

		b.putLine(strings.TrimRight(line, "\n\r"), b.curline)
		byCount += len(line)
	}

	if err == io.EOF {
		err = nil
	}

	return byCount, err
}

// Read fulfills io.Reader
func (b *scratchBuf) Read(p []byte) (int, error) {
	var byCount int
	var buf bytes.Buffer

	for idx := 1; idx <= b.getLastline(); idx++ {
		if byCount >= len(p) {
			return byCount, io.ErrShortBuffer
		}

		line := b.getLine(idx)
		buf.WriteString(line)
		buf.WriteRune('\n')
	}
	byCount = buf.Len()
	buf.Read(p)
	return byCount, io.EOF
}


func (b *scratchBuf) setCurline(idx int) {
	b.curline = idx
}

func (b *scratchBuf) getCurline() int {
	return b.curline
}

func (b *scratchBuf) setLastline(idx int) {
	b.lastline = idx
}

// getLastline reports how many lines are in the *active* buffer which is what got when we called getNumLines()
func (b *scratchBuf) getLastline() int {
	return b.lastline
}

// putLine adds a new lines to the end of the buffer then moves them into place
func (b *scratchBuf) putLine(line string, idx int) error {
	// NOTE: we are not guarding this index here
	b.curline = idx

	b.lastline++

	newLine, err := b.writeLine(line)
	if err != nil {
		return err
	}

	// when the current active buffer has fewer lines that the total buffer
	// (ie we've deleted/forgotten lines at the end of the buffer)
	// we can reuse those lines in stead of always appending.
	// This reduces memory usage but the scratch file will continue to grow.
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

	newLine, err := b.writeLine(line)
	if err != nil {
		return err
	}

	b.lines[idx] = newLine
	b.setDirty(true)
	return nil
}

// duplicateLine changes the line to the new text at the given index
func (b *scratchBuf) duplicateLine(idx int) error {
	if idx < 1 || idx > b.getLastline() {
		return fmt.Errorf("cannot duplicate text; invalid address; %d", idx)
	}

	b.lastline++ // regardless of where we put the duped line, the buffer is now one line greater

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
func (b *scratchBuf) nextLine(idx int) int {
	if idx >= b.lastline {
		return 0
	}
	return idx + 1
}

// prevLine returns the previous index in the buffer, looping to lastline after 0
func (b *scratchBuf) prevLine(idx int) int {
	if idx <= 0 {
		return b.lastline
	}
	return idx - 1
}

// returns the text of the line at the given index
func (b *scratchBuf) getLine(idx int) string {
	bts := make([]byte, b.lines[idx].len)
	_, err := b.ext.ReadAt(bts, int64(b.lines[idx].pos))
	if err != nil {
		stderr.Fatal(err)
	}

	return string(bts)
}

// scanForward returns a func that walks the buffer's indices in a forward loop.
// As an implementation detail, the number of lines is the number of non-zero lines.
func (b *scratchBuf) scanForward(start, num int) func() (int, bool) {
	stop := false
	idx := b.prevLine(start) // remove 1 because nextLine advances one

	if num < 0 {
		num = b.getLastline()
	}

	n := 0 // this is essentially a do{}while() loop so '0' will execute once
	return func() (int, bool) {
		if n <= num { // if we loop back, we'll return '0' so have to stop a '>' and not '>="
			n++
			idx = b.nextLine(idx)
			return idx, !stop
		}

		return idx, stop
	}
}

// scanReverse returns a func that walks the buffer's indices in a reverse loop
func (b *scratchBuf) scanReverse(start, num int) func() (int, bool) {
	stop := false
	idx := b.nextLine(start) // remove 1 because nextLine advances one

	if num < 0 {
		num = b.getLastline()
	}

	n := 0
	return func() (int, bool) {
		if n <= num { // if we loop back, we'll return '0' so have to stop a '>' and not '>="
			n++
			idx = b.prevLine(idx)
			return idx, !stop
		}

		return idx, stop
	}
}

func (b *scratchBuf) writeLine(line string) (bufferLine, error) {
	bts := strings.TrimRightFunc(line, func(r rune) bool {
		return unicode.IsSpace(r)
	})

	num, err := b.ext.WriteAt([]byte(bts), int64(b.pos))
	if err != nil {
		return bufferLine{}, err
	}

	newline := bufferLine{
		pos:  b.pos,
		len:  num,
		mark: null,
	}

	b.pos += num // track tha last bytes written because we'll start there next time

	return newline, nil
}

func (b *scratchBuf) destructor() {
	// if f, ok := b.ext.(*os.File); ok {
	// stderr.Println(f.Name())
	os.Remove(b.ext.Name())
	// }
}

func (b *scratchBuf) clone() buffer {
	cll := make([]bufferLine, len(b.lines))
	copy(cll, b.lines)
	t := &scratchBuf{
		curline:  b.curline,
		lastline: b.lastline,
		lines:    cll,
		filename: b.filename,
		ext:      b.ext,
		pos:      b.pos,
		dirty:    b.dirty,
		rev:      b.rev,
	}
	return t
}

func (b *scratchBuf) String() string {
	var rtn strings.Builder
	fmt.Fprint(&rtn, "buffer (file):\r\n")
	fmt.Fprintf(&rtn, "  filename: %s\r\n", b.filename)
	fmt.Fprintf(&rtn, "  curline: %d\r\n", b.curline)
	fmt.Fprintf(&rtn, "  lastline: %d\r\n", b.lastline)
	fmt.Fprintf(&rtn, "  dirty: %t\r\n", b.dirty)
	fmt.Fprintf(&rtn, "  rev: %d\r\n", b.rev)
	fmt.Fprintf(&rtn, "  pos: %d\r\n", b.pos)
	fmt.Fprintf(&rtn, "  scratch: %s\r\n", b.ext.Name())
	for k, v := range b.lines {
		if k == 0 {
			continue
		}
		fmt.Fprintf(&rtn, "  line[%d]: %s\r\n", k, v.String())
	}
	return rtn.String()
}

// defaultLines normalizes two addresses, both optional.
// It takes what is provided and returns sensible defaults with an eye to
// how the relate to each other.
// It also changes '.' and '$' to current and end addresses respectively
func (b *scratchBuf) defaultLines(num1, num2, incr string, line1, line2 int) (int, int, error) {
	var (
		err        error
		idx1, idx2 int
	)

	if len(num1) == 0 && len(num2) == 0 {
		return line1, line2, nil
	}

	idx1, err = b.makeAddress(num1, b.getCurline())
	if err != nil {
		return 0, 0, fmt.Errorf("invalid buffer start address: %s", err.Error())
	}

	idx2, err = b.makeAddress(num2, idx1)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid buffer end address: %s", err.Error())
	}

	if incr != "" {
		// this will coerce values that are out of range into range
		idx1, idx2 = b.applyIncrement(incr, idx1, idx2)
	}

	if idx1 > idx2 || idx1 <= 0 || idx1 > b.getLastline() || idx2 > b.getLastline() {
		return 0, 0, fmt.Errorf("invalid buffer range; %d, %d", idx1, idx2)
	}

	return idx1, idx2, nil
}

func (b *scratchBuf) applyIncrement(incr string, num1, relative int) (int, int) {
	num2 := relative
	if incr == ">" {
		num2 = num1 + relative
	}

	if incr == "<" {
		num2 = num1
		num1 -= relative
	}

	if num1 < 1 {
		num1 = 1
	}

	if num2 > b.getLastline() {
		num2 = b.getLastline()
	}

	return num1, num2

}

func (b *scratchBuf) makeContext(line1, line2, pager int) (int, int, error) {
	line1 = line1 - pager
	if line1 < 0 {
		line1 = 1
	}

	line2 = line2 + pager
	if !b.hasAddress(line2) {
		line2 = b.getLastline()
	}

	return line1, line2, nil
}

// converts a string address into a number with special cases for '.', '$', and ”.
// Start/end addresses are guarded against '0' elsewhere (in defaultLines) but
// are allowed in destinations
func (b *scratchBuf) makeAddress(addr string, def int) (int, error) {
	var (
		idx int
		err error
	)

	switch true {
	case addr == "$":
		idx = b.getLastline()
	case addr == ".":
		idx = b.getCurline()
	case addr == "":
		idx = def
	default:
		idx, err = intval(addr)
		if err != nil {
			return 0, fmt.Errorf("invalid address; %s", err.Error())
		}
	}

	// if idx < 0 || idx > b.getLastline() {
	// 	return 0, fmt.Errorf("unknown address: %s", addr)
	// }

	return idx, nil
}

func (b *scratchBuf) hasAddress(idx int) bool {
	if idx < 0 || idx > b.getLastline() {
		return false
	}
	return true
}

// delLines moves a range of lines to the end of the buffer
// then decreases the last line to "forget" about the lines at the end
func (b *scratchBuf) delLines(line1, line2 int) error {
	// NOTE: there is no validation here for accidentaly deleting all the lines; validation should be up a level
	if line1 <= 0 {
		line1 = 1
	}

	if line2 < line1 {
		return fmt.Errorf("unable to delete invalid range: %d,%d", line1, line2)
	}

	lastLine := b.getLastline()
	b.bulkMove(line1, line2, lastLine)
	b.setLastline(lastLine - (line2 - line1 + 1))
	b.setCurline(b.prevLine(line1))
	return nil
}
