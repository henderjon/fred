package main

import "errors"

type cache struct {
	pager       int
	column      int
	prevSearch  search
	prevReplace replace
	undos       []buffer
	undoPos     int
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

func (c *cache) setPreviousSearch(s search) {
	c.prevSearch = s
}

func (c *cache) getPreviousSearch() search {
	return c.prevSearch
}

func (c *cache) setPreviousReplace(s replace) {
	c.prevReplace = s
}

func (c *cache) getPreviousReplace() replace {
	return c.prevReplace
}

func (c *cache) stageUndo(b buffer) {
	c.undoPos++
	c.undos = append(c.undos, b)
}

func (c *cache) unstageUndo() (buffer, error) {
	if len(c.undos) > 0 {
		c.undoPos--
		return c.undos[c.undoPos], nil
	}
	return *(new(buffer)), errors.New("nothing to undo") // dereference the pointer before returning it
}

type stager interface {
	stageUndo(b buffer)
}
