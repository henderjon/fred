package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

// doPrintAddress asks for 'l2' because it should print the end of the requested range knowing that if only one address is given, it is a range of a single number
func doPrintAddress(b buffer, l2 int) error {
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
		cl = l3 + (l2 - l1 + 1) // the last line + the number of lines we moved (the difference of the origin range)
	}

	b.setCurline(cl)
	return nil
}

func doCopyNPaste(b buffer, l1, l2 int, dest string) error {
	var err error

	l3, err := guardAddress(dest, b.getCurline(), b.getLastline())
	if err != nil {
		return err
	}

	// flag where we're going to be adding lines
	mark := b.getLastline()

	// add old lines to the end of the buffer; we'll move them later
	b.setCurline(mark)

	for i := l1; i <= l2; i++ {
		err = b.putText(b.getText(i))
		if err != nil {
			return err
		}
	}

	mark++ // we added our new content here, let's move it to where we want it to be
	b.bulkMove(mark, mark+(l2-l1), l3)

	return nil
}

func doSimpleReplace(b buffer, l1, l2 int, pattern, replace, num string) error {
	var err error

	n, err := strconv.Atoi(num)
	if err != nil {
		return fmt.Errorf("invalid number: %s; %s", num, err.Error())
	}

	for idx := l1; idx <= l2; idx++ {
		old := b.getText(idx)
		new := strings.Replace(old, pattern, replace, n)
		b.replaceText(new, idx)
	}

	return err
}

// func doRegexReplace(b buffer, l1, l2 int, pattern, replace, num string) error {
// 	var (
// 		err error
// 		re  *regexp.Regexp
// 	)

// 	re, err = regexp.Compile(string(pat))
// 	if err != nil {
// 		return err
// 	}

// 	template := []byte(sub)

// 	result := []byte{}

// 	n := -1 // assume global; -1 no limits
// 	if !gflag {
// 		n = 1
// 	}

// 	for idx := addr.from; idx <= addr.to; idx++ {
// 		var (
// 			p   int
// 			old = getText(idx)
// 			new []byte
// 		)

// 		// go has no regex func for only doing ONE replacement. This uses a
// 		// workaround to walk through the string and manually replace each match
// 		// in order to emulate the behavior.

// 		// this finds the indexes of the matches
// 		submatches := re.FindAllSubmatchIndex(old, n)
// 		for n := 0; n < len(submatches); n++ {
// 			// expand any $1 replacements; this takes the text input 'old' and
// 			// using the indexes from 'submatches[n]' replaces it with the
// 			// expanded replacement in 'template' and appends it to 'result' in
// 			// other words, result is what should go into the new string
// 			result := re.Expand(result, template, old, submatches[n])
// 			// create a new string add the characters of the old string from the
// 			// beginning of the last match (or zero) to the beginning of the
// 			// current match (we're currently iterating)
// 			new = append(new, old[p:submatches[n][0]]...)
// 			// add the replacement value to the new string
// 			new = append(new, result...)
// 			// move the cursor to the index of the end of the current match so
// 			// that then we add from the index of the end of the current match to
// 			// the index of the beginning of the next match of the old string to
// 			// the new string. in effect, make sure we add the bits of the old
// 			// string that didn't match to the new string.
// 			p = submatches[n][1]
// 		}

// 		new = append(new, old[p:]...)
// 		replaceText(new, idx)
// 	}
// 	return err
// }
