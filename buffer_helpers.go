package main

import (
	"fmt"
	"strconv"
)

func makeContext(b buffer, l1, l2, pager int) (int, int, error) {
	l1 = l1 - pager
	if l1 < 0 {
		l1 = 1
	}

	l2 = l2 + pager
	if l2 > b.getLastline() {
		l2 = b.getLastline()
	}

	return l1, l2, nil
}

// converts a string address into a number with special cases for '.', '$', and ”. Start/end addresses are guarded against '0' elsewhere (in defaultLines) but allowed in destinations
func guardAddress(b buffer, addr string) (int, error) {
	if addr == "." || addr == "" { // if no address was given, use the current line
		return b.getCurline(), nil
	}

	if addr == "$" {
		return b.getLastline(), nil
	}

	i, err := strconv.Atoi(addr)
	if err != nil {
		return 0, fmt.Errorf("invalid address: %s; %s", addr, err.Error())
	}

	if i < 0 || i > b.getLastline() {
		return 0, fmt.Errorf("invalid address: %s", addr)
	}

	return i, nil
}
