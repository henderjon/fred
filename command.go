package main

import (
	"fmt"
	"strings"
)

type command struct {
	addrStart    string
	addrEnd      string
	addrIncr     string
	addrPattern  string
	action       rune
	pattern      string
	substitution string
	replaceNum   string
	destination  string
	subCommand   rune
	argument     string
	globalPrefix string
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
	fmt.Fprintf(&cmd, " replaceNum(%s)", c.replaceNum) // /g suffix replace nth/all match/es
	fmt.Fprintf(&cmd, " destination(%s)", c.destination)
	fmt.Fprintf(&cmd, " subCommand(%s)", string(c.subCommand))
	fmt.Fprintf(&cmd, " argument(%s)", c.argument)
	fmt.Fprintf(&cmd, " globalPrefix(%s)", c.globalPrefix) // g/ prefix; find more than one line

	return cmd.String()
}

func (c *command) setAddr(f string) {
	if len(c.addrIncr) > 0 ||
		len(c.addrStart) > 0 {
		c.addrEnd = f
		return
	}

	c.addrStart = f
}

func (c *command) setAction(a rune) {
	if c.action == 0 {
		c.action = rune(a)
	} else {
		c.setSubCommand(a)
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

func (c *command) setGlobalPrefix(s string) {
	c.globalPrefix = s
}

func (c *command) setReplaceNum(s string) {
	if s == string(globalReplaceAction) {
		c.replaceNum = "-1" // zero should mean all but is the zero value
		return
	}

	c.replaceNum = s
}

func (c *command) setSubCommand(r rune) {
	c.subCommand = r
}

func (c *command) setArgument(s string) {
	c.argument = s
}

func (c *command) setDestination(s string) {
	c.destination = s
}

func (c *command) setIncr(s string) {
	c.addrIncr = s
}
