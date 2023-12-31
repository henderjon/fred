package main

import (
	"fmt"
	"strings"
)

type cache struct {
	pager       int
	column      int
	prevSearch  search
	prevReplace replace
	prevBuffer  buffer
	currBuffer  buffer
	printType   int
}

func (c *cache) setPager(n int) {
	c.pager = n
}

func (c *cache) getPager() int {
	return c.pager
}

func (c *cache) setColumn(n int) {
	c.column = n
}

func (c *cache) getColumn() int {
	return c.column
}

func (c *cache) setPrintType(printType int) {
	c.printType = printType
}

func (c *cache) getPrintType() int {
	return c.printType
}

func (c *cache) search(pattern string, reverse bool) search {
	if len(pattern) == 0 { // no pattern means to repeat the last search
		c.prevSearch.reverse = reverse // allow the user to change directions
		return c.prevSearch
	}

	c.prevSearch = search{
		pattern: pattern,
		reverse: reverse,
	}

	return c.prevSearch
}

func (c *cache) replace(pattern, sub, num string) replace {
	if len(pattern) == 0 { // no pattern means to repeat the last search
		return c.prevReplace
	}

	c.prevReplace = replace{
		pattern:    pattern,
		substitute: sub,
		replaceNum: num,
	}
	return c.prevReplace
}

func (c *cache) getCurrBuffer() buffer {
	return c.currBuffer
}

func (c *cache) stageUndo(b buffer) {
	if c.currBuffer != nil && c.currBuffer.getRev() == b.getRev() {
		return
	}

	c.prevBuffer = c.currBuffer
	c.currBuffer = b
}

func (c *cache) unstageUndo() (buffer, error) {
	switch true {
	case c.prevBuffer != nil:
		tmp := c.prevBuffer
		c.stageUndo(c.prevBuffer)
		return tmp, nil
	case c.currBuffer != nil:
		return c.currBuffer, nil
	default:
		return *(new(buffer)), fmt.Errorf("nothing to undo") // dereference the pointer before returning it
	}
}

// String satisfies fmt.Stringer. It is useful for debugging and testing
// because it can expose the value of private properties and nested structs
func (c *cache) String() string {
	var rtn strings.Builder

	fmt.Fprint(&rtn, "cache:\r\n")

	if len(c.prevSearch.pattern) > 0 {
		fmt.Fprintf(&rtn, "  search.pattern: %s\r\n", c.prevSearch.pattern)
		fmt.Fprintf(&rtn, "  search.reverse: %t\r\n", c.prevSearch.reverse)
	}

	if len(c.prevReplace.pattern) > 0 {
		fmt.Fprintf(&rtn, "  replace.pattern: %s\r\n", c.prevReplace.pattern)
		fmt.Fprintf(&rtn, "  replace.replace: %s\r\n", c.prevReplace.substitute)
		fmt.Fprintf(&rtn, "  replace.replaceNum: %s\r\n", c.prevReplace.replaceNum)
	}

	fmt.Fprintf(&rtn, "  pager: %d\r\n", c.pager)
	fmt.Fprintf(&rtn, "  column: %d\r\n", c.column)

	if c.prevBuffer != nil {
		fmt.Fprintf(&rtn, "  prevBuffer: %d\r\n", c.prevBuffer.getRev())
	} else {
		fmt.Fprint(&rtn, "  prevBuffer: nil\r\n")
	}

	if c.currBuffer != nil {
		fmt.Fprintf(&rtn, "  currBuffer: %d\r\n", c.currBuffer.getRev())
	} else {
		fmt.Fprint(&rtn, "  currBuffer: nil\r\n")
	}

	return rtn.String()
}

// infinite undo but needs redo
// func (c *cache) stageUndo(b buffer) {
// 	c.undoPos++
// 	c.undos = append(c.undos, b)
// }

// func (c *cache) unstageUndo() (buffer, error) {
// 	if len(c.undos) > 0 {
// 		c.undoPos--
// 		return c.undos[c.undoPos], nil
// 	}
// 	return *(new(buffer)), errors.New("nothing to undo") // dereference the pointer before returning it
// }
