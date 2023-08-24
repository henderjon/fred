package main

type cache struct {
	pager       int
	column      int
	prevSearch  search
	prevReplace replace
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
