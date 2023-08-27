package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	printTypeReg = iota
	printTypeNum
	printTypeLit
	printTypeCol
)

func doPrint(inout termio, b buffer, l1, l2 int, cache *cache, printType int) error {
	var err error
	if l1 <= 0 || l1 > b.getNumLines() { // NOTE: l2 is not bound by last line; may be a problem
		return fmt.Errorf("unable to print; invalid address; %d; %d", l1, l2)
	}

	// b.setCurline(l1)

	l1, l2, err = makeContext(b, l1, l2, cache.getPager())
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

		if n > b.getNumLines() {
			break
		}

		line := b.getLine(n)
		switch printType {
		default:
			inout.printf("%s", line)
		case printTypeNum:
			inout.printf("%-2s%d\t%s", mk, n, line)
		case printTypeLit:
			inout.printf("%-2s%d\t%+q", mk, n, line)
		case printTypeCol:
			inout.printf("%-2s%d\t%s", mk, n, revealColumn(cache.getColumn(), line))
		}
	}

	b.setCurline(l2)

	return nil
}

func setPager(num string, cache *cache) (string, error) {
	if len(num) > 0 {
		var (
			err error
			n   int
		)

		n, err = intval(num)
		if err != nil {
			return "", fmt.Errorf("unable to set pager; %s", err.Error())
		}

		cache.setPager(n)
	}

	return fmt.Sprintf("pager set to %d", cache.getPager()), nil
}

// doPrintAddress asks for 'l2' because it should print the end of the requested range knowing that if only one address is given, it is a range of a single number
func doPrintAddress(b buffer, l2 int) (string, error) {
	b.setCurline(l2)
	return fmt.Sprintf("%d", b.getCurline()), nil
}

func doAppend(inout termio, b buffer, l1 int) error {
	return b.insertAfter(inout, l1)
}

func doInsert(inout termio, b buffer, l1 int) error {
	if l1 <= 1 {
		return b.insertAfter(inout, 0)
	}
	return b.insertAfter(inout, l1)
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

func doChange(inout termio, b buffer, l1, l2 int) error {
	err := doDelete(b, l1, l2)
	if err != nil {
		return err
	}
	return b.insertAfter(inout, b.prevLine(l1))
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

	for i := l1; i <= l2; i++ { // do this in reverse because we're putting them in order behind l3
		err = b.duplicateLine(i)
		if err != nil {
			return err
		}
	}

	mark++ // we added our new content here, let's move it to where we want it to be
	b.bulkMove(mark, mark+(l2-l1), l3)

	b.setCurline(l3)
	return nil
}

func doSimpleReplace(b buffer, l1, l2 int, pattern, sub, num string, cache *cache) error {
	var err error

	prevReplace := cache.getPreviousReplace()
	if len(pattern) == 0 { // no pattern means to repeat the last search
		pattern = prevReplace.pattern
		sub = prevReplace.replace
		num = prevReplace.replaceNum
	}

	cache.setPreviousReplace(replace{
		pattern:    pattern,
		replace:    sub,
		replaceNum: num,
	})

	n := 1 // default to 1; not -1 ("global")
	if len(num) > 0 {
		n, err = intval(num)
		if err != nil {
			return fmt.Errorf("unable to do a simple replace; %s", err.Error())
		}
	}

	for idx := l1; idx <= l2; idx++ {
		var new string
		old := b.getLine(idx)
		if n < 0 {
			// replace first n matches; always -1 (all)
			new = strings.Replace(old, pattern, sub, n)
		} else {
			// replace nth match
			new = simpleNReplace(old, pattern, sub, n)
		}
		b.replaceLine(new, idx)
	}

	b.setCurline(l2)
	return err
}

func doRegexReplace(b buffer, l1, l2 int, pattern, sub, num string, cache *cache) error {
	var (
		err error
		re  *regexp.Regexp
	)

	prevReplace := cache.getPreviousReplace()
	if len(pattern) == 0 { // no pattern means to repeat the last search
		pattern = prevReplace.pattern
		sub = prevReplace.replace
		num = prevReplace.replaceNum
	}

	cache.setPreviousReplace(replace{
		pattern:    pattern,
		replace:    sub,
		replaceNum: num,
	})

	pattern = handleTabs(pattern)
	// sub = handleTabs(prevReplace.replace)

	re, err = regexp.Compile(pattern)
	if err != nil {
		return err
	}

	result := []byte{}

	n := 1 // default to 1; not -1 ("global")
	if len(num) > 0 {
		n, err = intval(num)
		if err != nil {
			return fmt.Errorf("unable to do a regex replace; %s", err.Error())
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
			result := re.ExpandString(result, sub, old, submatches[i]) // submatches is [][]int .. the inner []int is the index of the beginning and the index of the end of the submatch
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

func doBreakLines(b buffer, l1, l2 int, pattern, sub, num string, cache *cache) error {
	if len(sub) > 0 {
		return doRegexReplace(b, l1, l2, pattern, sub, num, cache)
	}

	var (
		err error
		re  *regexp.Regexp
	)

	// our scan takes an upper bound number of iterations
	numLines := l2 - l1 // scan will handle <0

	err = doMarkLines(b, l1, numLines, pattern, false)
	if err != nil {
		return err
	}

	re, err = regexp.Compile(pattern)
	if err != nil {
		return err
	}

	n := 1 // default to 1; not -1 ("global")
	if len(num) > 0 {
		n, err = intval(num)
		if err != nil {
			return fmt.Errorf("unable to break lines; %s", err.Error())
		}
	}

	// scan backwards because adding lines to the buffer will screw up our scan
	scan := b.scanReverse(l2, numLines)
	for {
		idx, ok := scan()
		if !ok {
			break
		}

		b.setCurline(idx)

		var p int
		old := b.getLine(idx)

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
			// for documentation see doRegexReplace
			b.putLine(old[p:submatches[i][1]])
			p = submatches[i][1]
		}
		b.putLine(old[p:])
		doDelete(b, idx, idx)
	}

	b.setCurline(l2)
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
func doReadFile(b buffer, l1 int, filename string) (string, error) {
	var err error
	b.setCurline(l1)

	_, err = normalizeFilePath(b, filename)
	if err != nil {
		return "", err
	}

	f, err := os.Open(b.getFilename())
	if err != nil {
		return "", err
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

	return fmt.Sprint(numbyt), err
}

// doReadFile adds the contents of filename and adds them to the buffer after l1
func doWriteFile(inout termio, b buffer, l1, l2 int, filename string) (string, error) {
	var err error
	b.setCurline(l1)

	if len(b.getFilename()) <= 0 && len(filename) <= 0 {
		filename, err = inout.input("filename? ")
		if err != nil {
			return "", err
		}
		if len(filename) <= 0 {
			return "", fmt.Errorf("cannot write empty file name")
		}
	}

	_, err = normalizeFilePath(b, filename)
	if err != nil {
		return "", err
	}

	// f, err := os.Create(filename)
	f, err := os.OpenFile(b.getFilename(), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return "", err
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

	b.setDirty(false)
	return fmt.Sprint(numbyt), err
}

func doExternalShell(b buffer, l1, l2 int, command string) func(readFromBuffer bool, stdout io.Writer) (string, error) {
	return func(readFromBuffer bool, stdout io.Writer) (string, error) {
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
		return "!", err
	}
}

func doSetFilename(b buffer, filename string) (string, error) {
	_, err := normalizeFilePath(b, filename)
	if err != nil {
		return "", err
	}
	return b.getFilename(), nil
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

func doGetMarkedLine(inout termio, b buffer, mark string) error {
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
			inout.printf("%2d) %s", i, b.getLine(i))
			b.setCurline(i)
		}
	}

	return nil
}

func doGetNextMatchedLine(inout termio, b buffer, pattern string, forward bool, cache *cache) error {
	prevSearch := cache.getPreviousSearch()
	if len(pattern) == 0 { // no pattern means to repeat the last search
		pattern = prevSearch.pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	cache.setPreviousSearch(search{
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
			inout.printf("%2d) %s", i, b.getLine(i))
			b.setCurline(i)
			return nil
		}
	}
	return nil
}

func doSetColumn(num string, cache *cache) (string, error) {
	if len(num) > 0 {
		var (
			err error
			n   int
		)

		n, err = intval(num)
		if err != nil {
			return "", fmt.Errorf("unable to set column; %s", err.Error())
		}

		cache.setColumn(n)
	}

	return fmt.Sprintf("column set to %d", cache.getColumn()), nil
}

func doDebug(b buffer) (string, error) {
	// if !debug {
	// return "", errors.New("debugging is not enabled; did you mean to use `-debug`?")
	// }
	return b.String(), nil
}
