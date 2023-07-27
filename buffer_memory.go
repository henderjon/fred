package main

import (
	"errors"
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

// func (b memoryBuf) getCurline() int {
// 	return b.curline
// }

// func (b memoryBuf) setCurline(i int) {
// 	b.curline = i
// }

func (b memoryBuf) insertAfter(idx int, global bool) error {
	var err error
	if global {
		return errors.New("global operations not allowed")
	} else {
		b.curline = idx
		for stdin.Scan() {
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
	}
	return err
}

func (b memoryBuf) putText(line []byte) error {
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

func (b memoryBuf) defaultLines(start, end string) (int, int, error) {
	line1, err := b.defaultLine(start)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid first address: %s", err.Error())
	}

	line2, err := b.defaultLine(end)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid first address: %s", err.Error())
	}

	if line1 > line2 || line1 <= 0 {
		return 0, 0, errors.New("invalid range")
	}
	return line1, line2, nil // page 188
}

func (b memoryBuf) defaultLine(addr string) (int, error) {
	if addr == "." || addr == "" {
		return b.curline, nil
	}

	if addr == "$" {
		return b.lastline, nil
	}

	return strconv.Atoi(addr)
}
