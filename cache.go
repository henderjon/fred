package main

type cache struct {
	prevSearch  search
	prevReplace replace
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
