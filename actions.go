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

	b.setCurline(l2)

	return nil
}

// doPrintAdress asks for 'l2' because it should print the end of the requested range knowing that if only one address is given, it is a range of a single number
func doPrintAdress(b buffer, l2 int) error {
	b.setCurline(l2)
	fmt.Printf("%d\n", b.getCurline())
	return nil
}

func doAppend(b buffer, l1 int) error {
	return b.insertAfter(l1)
}

func doInsert(b buffer, l1 int) error {
	if l1 <= 1 {
		return b.insertAfter(0)
	}
	return b.insertAfter(l1)
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

func doChange(b buffer, l1, l2 int) error {
	err := doDelete(b, l1, l2)
	if err != nil {
		return err
	}
	return b.insertAfter(b.prevLine(l1))
}

func doMove(b buffer, l1, l2 int, dest string) error {
	l3, err := guardAddress(dest, b.getCurline(), b.getLastline())
	if err != nil {
		return err
	}

	// guard against bad addressing
	if (l1 <= 0 || l3 >= l1) && (l3 <= l2) {
		return fmt.Errorf("invalid ranges; move '%d' through '%d' to '%d'?", l1, l2, l3)
	}

	b.bulkMove(l1, l2, l3)
	var cl int
	if l3 > l1 {
		cl = l3
	} else {
		cl = l3 + (l2 - l1 + 1) // the last line + the difference of the origin range (should be a negative number)
	}

	b.setCurline(cl)
	return nil
}
