package main

import (
	"fmt"
	"strconv"
	"strings"
)

type command struct {
	addRange     []int
	addPattern   string
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

	fmt.Fprintf(&cmd, "%d", c.addRange[0])

	if len(c.addRange) >= 2 {
		fmt.Fprintf(&cmd, ",%d", c.addRange[1])
	}

	if len(c.addPattern) > 0 {
		fmt.Fprintf(&cmd, "; %s", c.addPattern)
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
	case c.numaddRange() == 0:
		if i < 0 {
			i = 0
		}
		c.addRange = append(c.addRange, i)
	case c.numaddRange() == 1:
		if i < 0 || i >= c.addRange[0] {
			c.addRange = append(c.addRange, i)
		} else {
			// repeat the first number if the second number is smaller than the first
			c.addRange = append(c.addRange, c.addRange[0])
			// TODO: use the $ end of the buffer for the last line
		}
	case c.numaddRange() >= 2: // TODO: maybe this should throw an error instead of compensating?
		if i < 0 || i >= c.addRange[0] {
			c.addRange[1] = i
		} else {
			c.addRange[1] = c.addRange[0]
		}
	}
}

func (c *command) numaddRange() int {
	return len(c.addRange)
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
