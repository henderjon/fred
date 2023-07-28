package main

import (
	"errors"
	"fmt"
)

func doPrint(b buffer, l1, l2 int) error {
	if l1 <= 0 || l1 > b.getNumLines() { // NOTE: l2 is not bound by last line; may be a problem
		return errors.New("doPrint; invalid address")
	}

	for n := l1; n <= l2; n++ {
		if n > b.getNumLines() {
			break
		}
		line := b.getText(n)
		fmt.Printf("%2d) %s\n", n, line)
	}

	b.setCurline(l1)

	return nil
}

func doAppend(b buffer, l1 int) error {
	return b.insertAfter(l1, false)
}

func doInsert(b buffer, l1 int) error {
	if l1 <= 1 {
		return b.insertAfter(0, false)
	}
	return b.insertAfter(l1, false)
}

// doDelete moves a range of lines to the end of the buffer then decreases the last line to "forget" about the lines at the end
func doDelete(b buffer, l1, l2 int) error {
	if l1 <= 0 {
		return errors.New("doDelete; invalid address")
	}

	ll := b.getLastline()
	b.bulkMove(l1, l2, ll)
	b.setLastline(ll - (l2 - l1 + 1))
	b.setCurline(b.prevLine(l1))
	return nil
}
