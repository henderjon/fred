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
