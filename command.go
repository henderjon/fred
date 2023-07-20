package main

import (
	"fmt"
	"strconv"
	"strings"
)

type command struct {
	addrs  []int
	action rune
	delim  rune
	regex  string
	sub    string
	// hasFrom  bool
	// hasTo    bool
	// hasDelim bool
}

func newCommand() *command {
	return &command{
		addrs:  make([]int, 0),
		action: noAction,
		// delim:  rune(unset),
		// regex:  "",
		// sub:    "",
	}
}

func (c *command) String() string {
	var cmd strings.Builder

	fmt.Fprintf(&cmd, "%d", c.addrs[0])

	if len(c.addrs) >= 2 {
		fmt.Fprintf(&cmd, ",%d", c.addrs[1])
	}

	fmt.Fprintf(&cmd, "%c%s%s%s%s%s",
		c.action,
		string(c.delim),
		c.regex,
		string(c.delim),
		c.sub,
		string(c.delim),
	)
	return cmd.String()
}

func (c *command) setFrom(f string) {
	i, e := strconv.Atoi(f)
	if e != nil {
		stderr.Log(e)
	}

	// as of now we only use the first and last numbers given,
	// to change this behavior to use only the first two numbers more the
	// `[1] = [0]` assignement to `case 1:` and drop all of `case 2:`
	switch c.numAddrs() {
	case 0:
		if i < 0 {
			i = 0
		}
		c.addrs = append(c.addrs, i)
	case 1:
		if i < 0 || i >= c.addrs[0] {
			c.addrs = append(c.addrs, i)
		} else {
			// repeat the first number if the second number is smaller than the first
			c.addrs = append(c.addrs, c.addrs[0])
			// 	// TODO use the $ end of the buffer for the last line
			// fmt.Fprintf(stderr, "%d, %d (%s)", c.addrs[0], i, logger.Here())
			// stderr.Log("second address must be larger than the first", )
		}
	case 2:
		if i < 0 || i >= c.addrs[0] {
			c.addrs[1] = i
		} else {
			c.addrs[1] = c.addrs[0]
		}
	default:
	}
}

func (c *command) numAddrs() int {
	return len(c.addrs)
}

// func (c *command) setFrom(f string) {
// 	i, e := strconv.Atoi(f)
// 	if e != nil {
// 		log.Println(e)
// 	}
// 	c.hasFrom = true
// 	c.from = i
// }

// func (c *command) setTo(t string) {
// 	i, e := strconv.Atoi(t)
// 	if e != nil {
// 		log.Println(e)
// 	}
// 	c.hasTo = true
// 	c.to = i
// }

func (c *command) setDelim(d string) {
	c.delim = rune(d[0])
}

func (c *command) setAction(a rune) {
	c.action = rune(a)
}
