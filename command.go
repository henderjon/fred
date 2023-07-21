package main

import (
	"fmt"
	"strconv"
	"strings"
)

type command struct {
	addrs        []int
	action       rune
	pattern      string
	substitution string
	additional   string
	// delim        rune
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

	fmt.Fprintf(&cmd, "; %c; %s; %s; %s",
		c.action,
		c.pattern,
		c.substitution,
		c.additional,
	)
	return cmd.String()
}

func (c *command) setAddr(f string) {
	i, e := strconv.Atoi(f)
	if e != nil {
		stderr.Log(e)
	}

	// as of now we only use the first and last numbers given,
	// to change this behavior to use only the first two numbers more the
	// `[1] = [0]` assignement to `case 1:` and drop all of `case 2:`
	switch true {
	default:
	case c.numAddrs() == 0:
		if i < 0 {
			i = 0
		}
		c.addrs = append(c.addrs, i)
	case c.numAddrs() == 1:
		if i < 0 || i >= c.addrs[0] {
			c.addrs = append(c.addrs, i)
		} else {
			// repeat the first number if the second number is smaller than the first
			c.addrs = append(c.addrs, c.addrs[0])
			// TODO: use the $ end of the buffer for the last line
		}
	case c.numAddrs() >= 2:
		if i < 0 || i >= c.addrs[0] {
			c.addrs[1] = i
		} else {
			c.addrs[1] = c.addrs[0]
		}
	}
}

func (c *command) numAddrs() int {
	return len(c.addrs)
}

func (c *command) setAction(a rune) {
	c.action = rune(a)
}

func (c *command) setPattern(s string) {
	c.pattern = s
}

func (c *command) setSubstitution(s string) {
	c.substitution = s
}
