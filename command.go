package main

import (
	"fmt"
	"strconv"
	"strings"
)

type command struct {
	addrStart    string
	addrEnd      string
	addrPattern  string
	action       rune
	pattern      string
	substitution string
	replaceNum   int
	destination  string
	subCommand   string
	globalPrefix bool
}

func (c *command) String() string {
	var cmd strings.Builder

	fmt.Fprintf(&cmd, "addrRange(%s", c.addrStart)

	if len(c.addrEnd) > 0 {
		fmt.Fprintf(&cmd, ",%s", c.addrEnd)
	}

	fmt.Fprintf(&cmd, ") addrPattern(%s)", c.addrPattern)
	fmt.Fprintf(&cmd, " action(%s)", string(c.action))
	fmt.Fprintf(&cmd, " pattern(%s)", c.pattern)
	fmt.Fprintf(&cmd, " substitution(%s)", c.substitution)
	fmt.Fprintf(&cmd, " replaceNum(%d)", c.replaceNum)
	fmt.Fprintf(&cmd, " destination(%s)", c.destination)
	fmt.Fprintf(&cmd, " subCommand(%s)", c.subCommand)
	fmt.Fprintf(&cmd, " globalPrefix(%t)", c.globalPrefix)

	return cmd.String()
}

func (c *command) setAddr(f string) {
	if len(c.addrStart) == 0 {
		c.addrStart = f
	} else {

		// if len(c.addrEnd) == 0 {
		c.addrEnd = f
	}
	return
	// var (
	// 	i   int
	// 	err error
	// )

	// if f == string(lastLine) { // handle the '$' special
	// 	i = -1
	// } else {
	// 	i, err = strconv.Atoi(f)
	// 	if err != nil {
	// 		stderr.Log(err)
	// 	}
	// }

	// // as of now we only use the first and last numbers given,
	// // to change this behavior to use only the first two numbers more the
	// // `[1] = [0]` assignment to `case 1:` and drop all of `case 2:`
	// switch true {
	// default:
	// case c.numaddrRange() == 0:
	// 	if i < 0 {
	// 		i = 0
	// 	}
	// 	c.addrRange = append(c.addrRange, i)
	// case c.numaddrRange() == 1:
	// 	if i < 0 || i >= c.addrRange[0] {
	// 		c.addrRange = append(c.addrRange, i)
	// 	} else {
	// 		// repeat the first number if the second number is smaller than the first
	// 		c.addrRange = append(c.addrRange, c.addrRange[0])
	// 		// TODO: use the $ end of the buffer for the last line
	// 	}
	// case c.numaddrRange() >= 2: // TODO: maybe this should throw an error instead of compensating?
	// 	if i < 0 || i >= c.addrRange[0] {
	// 		c.addrRange[1] = i
	// 	} else {
	// 		c.addrRange[1] = c.addrRange[0]
	// 	}
	// }
}

func (c *command) setAction(a rune) {
	if c.action == 0 {
		c.action = rune(a)
	} else {
		c.setSubCommand(string(a))
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

func (c *command) setReplaceNum(s string) {
	if s == string(gReplaceAction) {
		c.replaceNum = -1 // TODO: this really ought to be 0 but 0 is the default and we want to default to 1?
		return
	}

	i, e := strconv.Atoi(s)
	if e != nil {
		stderr.Log(e)
	}
	c.replaceNum = i
}

func (c *command) setSubCommand(s string) {
	c.subCommand = s
}

func (c *command) setDestination(s string) {
	c.destination = s
}
