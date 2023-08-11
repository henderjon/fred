package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
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

type stringable interface {
	string | rune | []byte
}

// contains checks haystack for needle, given that needle can be cast to a string; guards against empty string false positives
func contains[T stringable](haystack string, needle T) bool {
	find := string(needle)
	if len(find) == 0 {
		return false
	}

	return strings.Contains(haystack, find)
}

type interactor func(prompt string) (string, error)

// getInput creates a func that prints a message and takes
func getInput(in io.Reader, out io.Writer) interactor {
	stdin := bufio.NewScanner(os.Stdin)
	return interactor(func(prompt string) (string, error) {
		fmt.Fprint(out, prompt)
		stdin.Scan()
		return stdin.Text(), stdin.Err()
	})
}

func simpleNReplace(subject, pattern, replace string, n int) string {
	idx := strings.Index(subject, pattern)
	if idx == -1 {
		return subject
	}

	count := 1
	for ; count < n; count++ {
		i := strings.Index(subject[idx+len(pattern):], pattern)
		if i == -1 {
			break
		}
		idx += i + len(pattern)
	}

	if idx+len(pattern) > len(subject) || count < n {
		return subject
	}

	var rtn strings.Builder
	rtn.WriteString(subject[:idx])
	rtn.WriteString(replace)
	rtn.WriteString(subject[idx+len(pattern):])
	return rtn.String()
}

// clearBuffer blanks the current buffer. Long term this should probably be added to the buffer interface and handle checking for a dirty buffer ...
func clearBuffer(b buffer) error {
	if b.isDirty() {
		return errDirtyBuffer
	}

	// some commands require addresses
	line1, line2, err := b.defLines("", "", "", b.getCurline(), b.getCurline())
	if err != nil {
		return err
	}

	return doDelete(b, line1, line2)
}

func doMarkLines(b buffer, line1, numLines int, pattern string, invert bool) error {
	if len(pattern) == 0 {
		return fmt.Errorf("missing search pattern")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	scan := b.scanForward(line1, numLines)
	for {
		i, ok := scan()
		if !ok {
			break
		}

		if re.MatchString(b.getLine(i)) != invert {
			b.putMark(i, mark)
		} else if b.getMark(i) == mark {
			// previously marked lines that match should be ignored
			b.putMark(i, null)
		}

	}
	return nil
}
