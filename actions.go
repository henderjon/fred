package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	printTypeReg = iota
	printTypeNum
	printTypeLit
)

func doPrint(b buffer, l1, l2, pager int, printType int) error {
	var err error
	if l1 <= 0 || l1 > b.getNumLines() { // NOTE: l2 is not bound by last line; may be a problem
		return fmt.Errorf("unable to print; invalid address; %d; %d", l1, l2)
	}

	// b.setCurline(l1)

	l1, l2, err = makeContext(b, l1, l2, pager)
	if err != nil {
		return err
	}

	for n := l1; n <= l2; n++ {
		if n == 0 {
			continue // hide the '0' line
		}

		mk := ""
		if m := b.getMark(n); m != null && m != mark {
			mk += string(m)
		}

		// if n == b.getCurline() {
		// mk += "\u2192" // →
		// mk += "\u2022" // •
		// mk += "\u2588" // █
		// }

		if n > b.getNumLines() {
			break
		}

		// gutter := 0
		// for total := b.getNumLines(); total > 0; total /= 10 {
		// 	gutter++
		// }

		// gutterStr := fmt.Sprintf("%%-2s%%%dd \u2502", gutter)
		// gutterStr = fmt.Sprintf(gutterStr, mark, n)

		line := b.getLine(n)
		switch printType {
		default:
			fmt.Printf("%-2s%s\n", mk, line)
		case printTypeNum:
			fmt.Printf("%-2s%d\t%s\n", mk, n, line)
		case printTypeLit:
			fmt.Printf("%-2s %d\t%+q\n", mk, n, line)
		}
	}

	b.setCurline(l2)

	return nil
}

func setPager(p *int, num string) error {
	if len(num) > 0 {
		var (
			err error
			n   int
		)

		n, err = strconv.Atoi(num)
		if err != nil {
			return fmt.Errorf("unable to set pager; invalid number: %s; %s", num, err.Error())
		}

		*p = n
	}

	fmt.Printf("pager set to %d\n", *p)
	return nil
}

// doPrintAddress asks for 'l2' because it should print the end of the requested range knowing that if only one address is given, it is a range of a single number
func doPrintAddress(b buffer, l2 int) error {
	b.setCurline(l2)
	fmt.Printf("%d\n", b.getCurline())
	return nil
}

func doAppend(input interactor, b buffer, l1 int) error {
	return b.insertAfter(input, l1)
}

func doInsert(input interactor, b buffer, l1 int) error {
	if l1 <= 1 {
		return b.insertAfter(input, 0)
	}
	return b.insertAfter(input, l1)
}

// doDelete moves a range of lines to the end of the buffer then decreases the last line to "forget" about the lines at the end
func doDelete(b buffer, l1, l2 int) error {
	if l1 <= 0 {
		l1 = 1
		// return fmt.Errorf("unable to delete; invalid addresses; %d, %d", l1, l2)
	}

	ll := b.getLastline()
	b.bulkMove(l1, l2, ll)
	b.setLastline(ll - (l2 - l1 + 1))
	b.setCurline(b.prevLine(l1))
	return nil
}

func doChange(input interactor, b buffer, l1, l2 int) error {
	err := doDelete(b, l1, l2)
	if err != nil {
		return err
	}
	return b.insertAfter(input, b.prevLine(l1))
}

func doMove(b buffer, l1, l2 int, dest string) error {
	l3, err := guardAddress(b, dest)
	if err != nil {
		return err
	}

	// guard against bad addressing
	if (l1 <= 0 || l3 >= l1) && (l3 <= l2) {
		return fmt.Errorf("invalid ranges; unable to move '%d' through '%d' to '%d'?", l1, l2, l3)
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

	l3, err := guardAddress(b, dest)
	if err != nil {
		return err
	}

	// flag where we're going to be adding lines
	mark := b.getLastline()

	// add old lines to the end of the buffer; we'll move them later
	b.setCurline(mark)

	for i := l1; i <= l2; i++ {
		err = b.putLine(b.getLine(i))
		if err != nil {
			return err
		}
	}

	mark++ // we added our new content here, let's move it to where we want it to be
	b.bulkMove(mark, mark+(l2-l1), l3)

	b.setCurline(l3)
	return nil
}

func doSimpleReplace(b buffer, l1, l2 int, pattern, replace, num string) error {
	var err error

	n := 1 // default to 1; not -1 ("global")
	if len(num) > 0 {
		n, err = strconv.Atoi(num)
		if err != nil {
			return fmt.Errorf("unable to do a simple replace; invalid number: %s; %s", num, err.Error())
		}
	}

	for idx := l1; idx <= l2; idx++ {
		var new string
		old := b.getLine(idx)
		if n < 0 {
			// replace first n matches; always -1 (all)
			new = strings.Replace(old, pattern, replace, n)
		} else {
			// replace nth match
			new = simpleNReplace(old, pattern, replace, n)
		}
		b.replaceLine(new, idx)
	}

	b.setCurline(l2)
	return err
}

func doRegexReplace(b buffer, l1, l2 int, pattern, replace, num string) error {
	var (
		err error
		re  *regexp.Regexp
	)

	re, err = regexp.Compile(pattern)
	if err != nil {
		return err
	}

	result := []byte{}

	n := 1 // default to 1; not -1 ("global")
	if len(num) > 0 {
		n, err = strconv.Atoi(num)
		if err != nil {
			return fmt.Errorf("unable to do a regex replace; invalid number: %s; %s", num, err.Error())
		}
	}

	for idx := l1; idx <= l2; idx++ {
		var (
			p   int
			old = b.getLine(idx)
			new strings.Builder
		)

		// go has no regex func for only doing ONE replacement. This uses a
		// workaround to walk through the string and manually replace each match
		// in order to emulate the behavior.

		// this finds the indexes of the matches
		submatches := re.FindAllStringSubmatchIndex(old, n)
		for i := 0; i < len(submatches); i++ {
			// this catch allows us to do nth replacements
			if n != -1 && n-1 != i { // adjust for 0 index
				continue
			}
			// expand any $1 replacements; this takes the text input 'old' and
			// using the indexes from 'submatches[n]' replaces it with the
			// expanded replacement in 'replace' and appends it to 'result' in
			// other words, result is what should go into the new string
			result := re.ExpandString(result, replace, old, submatches[i])
			// create a new string add the characters of the old string from the
			// beginning of the last match (or zero) to the beginning of the
			// current match (we're currently iterating)
			new.WriteString(old[p:submatches[i][0]])
			// add the replacement value to the new string
			new.WriteString(string(result))
			// move the cursor to the index of the end of the current match so
			// that then we add from the index of the end of the current match to
			// the index of the beginning of the next match of the old string to
			// the new string. in effect, make sure we add the bits of the old
			// string that didn't match to the new string.
			p = submatches[i][1]
		}

		new.WriteString(old[p:])
		b.replaceLine(new.String(), idx)
		b.setCurline(idx)
	}
	return err
}

func doJoinLines(b buffer, l1, l2 int, sep string) error {
	var (
		err error
		new strings.Builder
	)

	if len(sep) == 0 {
		sep = " "
	}

	// this should prevent putLine() from moving the lines since we'll be doing
	// it ourselves
	b.setCurline(b.getLastline())
	for idx := l1; idx <= l2; idx++ {
		// copy all lines and combine them
		old := b.getLine(idx)
		new.WriteString(strings.TrimSpace(old))
		new.WriteString(sep)
	}

	// add them to the end of the bufferLines
	err = b.putLine(strings.TrimSuffix(new.String(), sep))
	if err != nil {
		return err
	}

	// then move them into place
	b.bulkMove(b.getCurline(), b.getCurline(), b.prevLine(l1))
	doDelete(b, b.nextLine(l1), b.nextLine(l2))
	b.setCurline(l1)
	return err
}

func doTransliterate(b buffer, l1, l2 int, pattern, replace string) error {
	var err error

	if utf8.RuneCountInString(pattern) != utf8.RuneCountInString(replace) {
		return fmt.Errorf("cannot transliterate; match and replace strings are different lengths")
	}

	replacements := make([]rune, 0)
	for _, r := range replace {
		replacements = append(replacements, r)
	}

	for idx := l1; idx <= l2; idx++ {
		var new strings.Builder
		old := b.getLine(idx)
		for _, oldRune := range old {
			newRune := oldRune
			for patternIdx, patternRune := range pattern {
				if oldRune == patternRune {
					if patternIdx < len(replacements) {
						newRune = replacements[patternIdx]
					}
				}
			}
			new.WriteRune(newRune)
		}
		b.replaceLine(new.String(), idx)
	}
	return err
}

func doMirrorLines(b buffer, l1, l2 int) error {
	b.reverse(l1, l2)
	return nil
}

// doReadFile adds the contents of filename and adds them to the buffer after l1
func doReadFile(b buffer, l1 int, filename string) error {
	var err error
	b.setCurline(l1)

	if len(filename) > 0 {
		b.setFilename(filename)
	}

	absPath, err := filepath.Abs(b.getFilename())
	if err != nil {
		return err
	}

	f, err := os.Open(absPath)
	if err != nil {
		return err
	}

	numbyt := 0
	fs := bufio.NewScanner(f)
	fs.Split(bufio.ScanLines)
	for fs.Scan() {
		err = fs.Err()
		if err != nil {
			break
		}

		numbyt += len(fs.Text()) + 1 // \n is always 1
		b.putLine(fs.Text())
	}

	fmt.Fprintln(os.Stdout, numbyt)
	return err
}

// doReadFile adds the contents of filename and adds them to the buffer after l1
func doWriteFile(b buffer, l1, l2 int, filename string) error {
	var err error
	b.setCurline(l1)

	if len(filename) > 0 {
		b.setFilename(filename)
	}

	// f, err := os.Create(filename)
	f, err := os.OpenFile(b.getFilename(), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	// TODO: only write l1 thru l2
	// for i := l1; i <= l2; i++ {
	// 	f.Write([]byte(b.getLine(i)))
	// 	f.Write([]byte{'\n'})
	// }

	numbyt := 0
	for i := 1; i <= b.getLastline(); i++ {
		n, _ := f.Write([]byte(b.getLine(i)))
		f.Write([]byte{'\n'})
		numbyt += n + 1 // \n is always 1
	}

	fmt.Fprintln(os.Stdout, numbyt)
	return err
}

func doExternalShell(b buffer, l1, l2 int, command string) func(readFromBuffer bool, stdout io.Writer) error {
	return func(readFromBuffer bool, stdout io.Writer) error {
		var (
			err   error
			stdin io.ReadWriter = nil
		)

		// fill a temp buffer to act as stdin
		if readFromBuffer {
			stdin = &bytes.Buffer{}
			for i := l1; i <= l2; i++ {
				stdin.Write([]byte(b.getLine(i)))
				stdin.Write([]byte("\n"))
			}
		}

		cmds := strings.TrimSpace(command)

		// hide all escaped '%'
		cmds = strings.ReplaceAll(cmds, `\%`, string(rune(26)))
		// replace '%' with filename
		cmds = strings.ReplaceAll(cmds, `%`, b.getFilename())
		// put the '%' back
		cmds = strings.ReplaceAll(cmds, string(rune(26)), `%`)

		buf := bufio.NewScanner(strings.NewReader(cmds))
		buf.Split(bufio.ScanWords)

		var args []string
		for buf.Scan() {
			args = append(args, os.ExpandEnv(buf.Text()))
		}

		cmd := exec.Command(args[0], args[1:]...)

		cmd.Stdin = stdin
		cmd.Stdout = stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		fmt.Fprintln(os.Stdout, "!")
		return err
	}
}

func doSetFilename(b buffer, filename string) error {
	if len(filename) > 0 {
		b.setFilename(filename)
	}

	fmt.Fprintln(os.Stdout, b.getFilename())

	return nil
}

func doSetMarkLine(b buffer, l1, l2 int, mark string) error {
	mk := null
	if len(mark) > 0 {
		mk = rune(mark[0])
	}

	for i := l1; i <= l2; i++ {
		b.putMark(i, mk)
	}
	return nil
}

func doGetMarkedLine(b buffer, mark string) error {
	mk := null
	if len(mark) > 0 {
		mk = rune(mark[0])
	}

	scan := b.scanForward(b.nextLine(b.getCurline()), b.getNumLines())
	for {
		i, ok := scan()
		if !ok {
			break
		}

		if b.hasMark(i, mk) {
			fmt.Printf("%2d) %s\n", i, b.getLine(i))
			b.setCurline(i)
		}
	}

	return nil
}

func doGetNextMatchedLine(b buffer, pattern string, forward bool) error {
	prevSearch := b.getPreviousSearch()
	if len(pattern) == 0 { // no pattern means to repeat the last search
		pattern = prevSearch.pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	b.setPreviousSearch(search{
		reverse: !forward,
		pattern: pattern,
	})

	scan := b.scanForward(b.nextLine(b.getCurline()), b.getNumLines())
	if !forward {
		scan = b.scanReverse(b.prevLine(b.getCurline()), b.getNumLines())
	}

	for {
		i, ok := scan()
		if !ok {
			break
		}

		if re.MatchString(b.getLine(i)) {
			fmt.Printf("%2d) %s\n", i, b.getLine(i))
			b.setCurline(i)
			return nil
		}
	}
	return nil
}
