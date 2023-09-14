package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

func makeContext(b buffer, l1, l2, pager int) (int, int, error) {
	l1 = l1 - pager
	if l1 < 0 {
		l1 = 1
	}

	l2 = l2 + pager
	if !b.hasAddress(l2) {
		l2 = b.getLastline()
	}

	return l1, l2, nil
}

type stringable interface {
	string | rune | []byte
}

// contains checks haystack for needle, given that needle can be cast to a string;
// guards against empty string false positives
func contains[T stringable](haystack string, needle T) bool {
	find := string(needle)
	if len(find) == 0 {
		return false
	}

	return strings.Contains(haystack, find)
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

// clearBuffer blanks the current buffer.
// Long term this should probably be added to the buffer interface and
// handle checking for a dirty buffer ...
func clearBuffer(b buffer) error {
	if b.isDirty() {
		return errDirtyBuffer
	}

	// some commands require addresses
	line1, line2, err := b.defaultLines("", "", "", 1, b.getLastline())
	if err != nil {
		return err
	}

	return doDelete(b, line1, line2)
}

func handleTabs(s string) string {
	s = strings.ReplaceAll(s, `\\t`, "\x1A") // \x1A is just a placeholder
	s = strings.ReplaceAll(s, `\t`, "\x09")
	return strings.ReplaceAll(s, "\x1A", `\t`)
}

// func handleNewlines(s string) string {
// 	s = strings.ReplaceAll(s, `\\n`, "\x1A")
// 	s = strings.ReplaceAll(s, `\n`, "\n")
// 	return strings.ReplaceAll(s, "\x1A", `\n`)
// }

func normalizeFilePath(b buffer, filename string) (string, error) {
	if len(filename) == 0 {
		// NOTE: there a far bit of paranoia here, this is very unlikely to happen IRL
		if len(b.getFilename()) == 0 {
			return "", errEmptyFilename
		}
		filename = b.getFilename()
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	b.setFilename(absPath)
	return absPath, nil
}

func revealColumn(col int, s string) string {
	if col <= 0 {
		return s
	}

	var rtn strings.Builder
	for i, r := range s {
		if i > 0 && ((i+1)%col) == 0 {
			rtn.WriteString("\033[7m")
			rtn.WriteRune(r)
			rtn.WriteString("\033[0m")
			continue
		}
		rtn.WriteRune(r)
	}
	return rtn.String()
}

func intval(num string) (int, error) {
	if len(num) == 0 {
		return 0, errors.New("empty number")
	}

	// if users are allowed to use a space between the action and the number, strip them first
	i, err := strconv.Atoi(strings.TrimSpace(num))
	if err != nil {
		return 0, fmt.Errorf("unable to parse number: %s; %s", num, err.Error())
	}
	return i, nil
}

func firstRune(s string) (rune, error) {
	if s == "" {
		return -1, fmt.Errorf("cannot parse rune: %s", s)
	}

	rn, wid := utf8.DecodeRuneInString(s)
	if wid == 0 {
		return -1, fmt.Errorf("cannot parse rune: %s", s)
	}
	return rn, nil
}
