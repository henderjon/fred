package main

import (
	"fmt"
	"strconv"
	"strings"
)

type command struct {
	addrRange    []int
	addrPattern  string
	action       rune
	pattern      string
	substitution string
	destination  int
	additional   string
	globalPrefix bool
	globalSuffix bool
}

func (c *command) String() string {
	var cmd strings.Builder

	fmt.Fprintf(&cmd, "%d", c.addrRange[0])

	if len(c.addrRange) >= 2 {
		fmt.Fprintf(&cmd, ",%d", c.addrRange[1])
	}

	if len(c.addrPattern) > 0 {
		fmt.Fprintf(&cmd, "; %s", c.addrPattern)
	}

	fmt.Fprintf(&cmd, "; %c; %d; %s; %s; %s",
		c.action,
		c.destination,
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
	case c.numaddrRange() == 0:
		if i < 0 {
			i = 0
		}
		c.addrRange = append(c.addrRange, i)
	case c.numaddrRange() == 1:
		if i < 0 || i >= c.addrRange[0] {
			c.addrRange = append(c.addrRange, i)
		} else {
			// repeat the first number if the second number is smaller than the first
			c.addrRange = append(c.addrRange, c.addrRange[0])
			// TODO: use the $ end of the buffer for the last line
		}
	case c.numaddrRange() >= 2: // TODO: maybe this should throw an error instead of compensating?
		if i < 0 || i >= c.addrRange[0] {
			c.addrRange[1] = i
		} else {
			c.addrRange[1] = c.addrRange[0]
		}
	}
}

func (c *command) numaddrRange() int {
	return len(c.addrRange)
}

func (c *command) setAction(a rune) {
	if c.action == 0 {
		c.action = rune(a)
	} else {
		c.setAdditional(string(a))
	}
}

func (c *command) setPattern(s string) {
	c.pattern = s
}

func (c *command) setAddrPattern(s string) {
	c.addrPattern = s
}

func (c *command) setSubstitution(s string) {
	c.substitution = s
}

func (c *command) setGlobalPrefix(b bool) {
	c.globalPrefix = b
}

func (c *command) setGlobalSuffix(b bool) {
	c.globalSuffix = b
}

func (c *command) setAdditional(s string) {
	c.additional = s
}

func (c *command) setDestination(s string) {
	i, e := strconv.Atoi(s)
	if e != nil {
		stderr.Log(e)
	}
	c.destination = i
}
