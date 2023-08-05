package main

import (
	"fmt"
	"strings"
)

type command struct {
	addrStart    string
	addrEnd      string
	addrPattern  string
	action       rune
	pattern      string
	substitution string
	replaceNum   string
	destination  string
	subCommand   string
	argument     string
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
	fmt.Fprintf(&cmd, " replaceNum(%s)", c.replaceNum) // /g suffix replace nth/all match/es
	fmt.Fprintf(&cmd, " destination(%s)", c.destination)
	fmt.Fprintf(&cmd, " subCommand(%s)", c.subCommand)
	fmt.Fprintf(&cmd, " argument(%s)", c.argument)
	fmt.Fprintf(&cmd, " globalPrefix(%t)", c.globalPrefix) // g/ prefix; find more than one line

	return cmd.String()
}

func (c *command) setAddr(f string) {
	if len(c.addrStart) == 0 {
		c.addrStart = f
		return
	}

	c.addrEnd = f
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
	if s == string(globalReplaceAction) {
		c.replaceNum = "-1" // zero means all
		return
	}

	c.replaceNum = s
}

func (c *command) setSubCommand(s string) {
	c.subCommand = s
}

func (c *command) setArgument(s string) {
	c.argument = s
}

func (c *command) setDestination(s string) {
	c.destination = s
}
