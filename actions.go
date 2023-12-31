package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	printTypeReg = iota
	printTypeNum
	printTypeLit
)

func doPrint(out io.Writer, b buffer, l1, l2 int, cache *cache, printType int) error {
	var err error

	// b.setCurline(l1)

	l1, l2, err = b.makeContext(l1, l2, cache.getPager())
	if err != nil {
		return err
	}

	// because these values were vetted by defaultLines there is no danger in
	// skipping validation within the scan
	scan := b.scanForward(l1, l2-l1)
	for {
		i, ok := scan() // scan hides the '0' line
		if !ok {
			break
		}

		mk := ""
		if m := b.getMark(i); m != null && m != mark {
			mk += string(m)
		}

		line := b.getLine(i)
		switch printType {
		default:
			fallthrough
		case printTypeNum:
			fmt.Fprintf(out, "%-2s%d\t%s", mk, i, line)
		case printTypeReg:
			fmt.Fprintf(out, "%s", line)
		case printTypeLit:
			fmt.Fprintf(out, "%-2s%d\t%+q", mk, i, line)
		}
	}

	cache.setPrintType(printType)
	b.setCurline(l2)

	return nil
}

func doSetPager(num string, cache *cache) (string, error) {
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

// doPrintAddress asks for 'l2' because it should print the end of the requested range
// knowing that if only one address is given, it is a range of a single number
func doPrintAddress(b buffer, l2 int) (string, error) {
	b.setCurline(l2)
	return fmt.Sprintf("%d", b.getCurline()), nil
}

func doAppend(inout termio, b buffer, l1 int) error {
	var line string
	var err error

	for {
		line, err = inout.input("")
		if err != nil {
			return err
		}

		if len(line) == 1 && line[0] == '.' {
			return nil
		}

		err = b.putLine(line, l1)
		if err != nil {
			return err
		}
		// append to the next line;
		// don't lock append to always adding lines to the one given,
		// move the destination with what is entered
		l1++
		b.setCurline(l1)
	}
}

func doInsert(inout termio, b buffer, l1 int) error {
	l1 = b.prevLine(l1)
	if l1 < 0 {
		l1 = 0
	}

	return doAppend(inout, b, l1)
}

// doDelete moves a range of lines to the end of the buffer
// then decreases the last line to "forget" about the lines at the end
func doDelete(b buffer, l1, l2 int) error {
	return b.delLines(l1, l2)
}

func doChange(inout termio, b buffer, l1, l2 int) error {
	err := doDelete(b, l1, l2)
	if err != nil {
		return err
	}
	return doAppend(inout, b, b.prevLine(l1))
}

func doMove(b buffer, l1, l2 int, dest string) error {
	l3, err := b.makeAddress(dest, b.nextLine(b.getCurline()))
	if err != nil {
		return err
	}

	// guard against bad addressing
	if !b.hasAddress(l3) {
		return fmt.Errorf("unable to move; destination our of range: %d", l3)
	}

	// guard against bad addressing
	if (l1 <= 0 || l3 >= l1) && (l3 <= l2) {
		return fmt.Errorf("invalid ranges; unable to move '%d' through '%d' to '%d'", l1, l2, l3)
	}

	b.bulkMove(l1, l2, l3)
	var cl int
	if l3 > l1 {
		cl = l3
	} else {
		// the last line + the number of lines we moved
		// ... the difference of the origin range
		cl = l3 + (l2 - l1 + 1)
	}

	b.setCurline(cl)
	return nil
}

func doCopyNPaste(b buffer, l1, l2 int, dest string) error {
	var err error

	l3, err := b.makeAddress(dest, b.nextLine(b.getCurline()))
	if err != nil {
		return err
	}

	// guard against bad addressing
	if !b.hasAddress(l3) {
		return fmt.Errorf("unable to paste; destination our of range: %d", l3)
	}

	// flag where we're going to be adding lines
	mark := b.getLastline()

	// add old lines to the end of the buffer; we'll move them later
	b.setCurline(mark)

	for idx := l1; idx <= l2; idx++ { // do this in reverse because we're putting them in order behind l3
		err = b.duplicateLine(idx)
		if err != nil {
			return err
		}
	}

	mark++ // we added our new content here, let's move it to where we want it to be
	b.bulkMove(mark, mark+(l2-l1), l3)

	b.setCurline(l3)
	return nil
}

func doSimpleReplace(b buffer, l1, l2 int, rep replace) error {
	var err error

	if rep.pattern == "" {
		return fmt.Errorf("empty pattern")
	}

	n := 1 // default to 1; not -1 ("global")
	if len(rep.replaceNum) > 0 {
		n, err = intval(rep.replaceNum)
		if err != nil {
			return fmt.Errorf("unable to do a simple replace; %s", err.Error())
		}
	}

	for idx := l1; idx <= l2; idx++ {
		var new string
		old := b.getLine(idx)
		if n < 0 {
			// replace first n matches; always -1 (all)
			new = strings.Replace(old, rep.pattern, rep.substitute, n)
		} else {
			// replace nth match
			new = simpleNReplace(old, rep.pattern, rep.substitute, n)
		}
		b.replaceLine(new, idx)
	}

	b.setCurline(l2)
	return err
}

func doRegexReplace(b buffer, l1, l2 int, rep replace) error {
	var (
		err error
		re  *regexp.Regexp
	)

	if rep.pattern == "" {
		return fmt.Errorf("empty pattern")
	}

	pattern := handleTabs(rep.pattern)
	// sub = handleTabs(prevReplace.replace)

	re, err = regexp.Compile(pattern)
	if err != nil {
		return err
	}

	result := []byte{}

	n := 1 // default to 1; not -1 ("global")
	if len(rep.replaceNum) > 0 {
		n, err = intval(rep.replaceNum)
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
			result := re.ExpandString(result, rep.substitute, old, submatches[i]) // submatches is [][]int .. the inner []int is the index of the beginning and the index of the end of the submatch
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

	// this should prevent putLine() from moving the lines since we'll be doing it ourselves
	b.setCurline(b.getLastline())
	for idx := l1; idx <= l2; idx++ {
		// copy all lines and combine them
		old := b.getLine(idx)
		new.WriteString(strings.TrimSpace(old))
		new.WriteString(sep)
	}

	// add them to the end of the bufferLines
	err = b.putLine(strings.TrimSuffix(new.String(), sep), b.getCurline())
	if err != nil {
		return err
	}

	// then move them into place
	b.bulkMove(b.getCurline(), b.getCurline(), b.prevLine(l1))
	doDelete(b, b.nextLine(l1), b.nextLine(l2))
	b.setCurline(l1)
	return err
}

func doBreakLines(b buffer, l1, l2 int, rep replace) error {
	if len(rep.substitute) > 0 {
		// if we're "breaking" by injecting, just proxy to regex replace
		return doRegexReplace(b, l1, l2, rep)
	}

	var (
		err error
		re  *regexp.Regexp
	)

	// our scan takes an upper bound number of iterations
	numLines := l2 - l1 // scan will handle <0

	// err = doMarkLinesRegex(b, l1, numLines, rep.pattern, false)
	// if err != nil {
	// 	return err
	// }

	re, err = regexp.Compile(rep.pattern)
	if err != nil {
		return err
	}

	n := 1 // default to 1; not -1 ("global")
	if len(rep.replaceNum) > 0 {
		n, err = intval(rep.replaceNum)
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
			b.putLine(old[p:submatches[i][1]], b.getCurline())
			p = submatches[i][1]
		}
		b.putLine(old[p:], b.getCurline())
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
	b.setCurline(l2)
	return nil
}

func doEditFile(b buffer, fs FileSystem, filename string) (string, error) {
	nb, err := doReadFile(b, 1, fs, filename)
	if err != nil {
		return nb, err
	}

	path, err := doSetFilename(b, fs, filename)
	if err != nil {
		return path, err
	}

	return fmt.Sprintf("%s; %s", nb, path), nil
}

// doReadFile adds the contents of filename and adds them to the buffer after l1
func doReadFile(b buffer, l1 int, fs FileSystem, filename string) (string, error) {
	var err error
	b.setCurline(l1)

	if len(filename) == 0 {
		filename = b.getFilename()
	}

	frdr, err := fs.FileReader(filename)
	if err != nil {
		return "", err
	}

	defer frdr.Close()
	b.setFilename(filename)

	numbyt := 0
	fscan := bufio.NewScanner(frdr)
	fscan.Split(bufio.ScanLines)

	for fscan.Scan() {
		err = fscan.Err()
		if err != nil {
			break
		}

		numbyt += len(fscan.Text()) + 1 // \n is always 1
		b.putLine(fscan.Text(), b.getCurline())
	}

	return fmt.Sprint(numbyt), err
}

// doReadFile adds the contents of filename and adds them to the buffer after l1
func doWriteFile(inout termio, b buffer, l1, l2 int, fs FileSystem, filename string) (string, error) {
	var err error
	b.setCurline(l1)

	if len(filename) == 0 {
		filename = b.getFilename()
	}

	if len(filename) <= 0 {
		filename, err = inout.input("filename? ")
		if err != nil {
			return "", err
		}
	}

	// TODO: only write l1 thru l2
	// for i := l1; i <= l2; i++ {
	// 	f.Write([]byte(b.getLine(i)))
	// 	f.Write([]byte{'\n'})
	// }

	fwrt, err := fs.FileWriter(filename)
	if err != nil {
		return "", err
	}

	defer fwrt.Close()

	var numbts int
	scan := b.scanForward(1, b.getLastline()-1) // don't loop all the way
	for {
		idx, ok := scan()
		if !ok {
			break
		}

		n, _ := fwrt.Write([]byte(b.getLine(idx)))
		fwrt.Write([]byte("\n"))
		numbts += n + 1
	}

	b.setDirty(false)
	return fmt.Sprint(numbts), err
}

func doExternalShell(b buffer, l1, l2 int, command string) func(stdin io.Reader, stdout io.Writer) (string, error) {
	return func(stdin io.Reader, stdout io.Writer) (string, error) {
		var (
			err error
			// stdin io.ReadWriter = nil
		)

		// fill a temp buffer to act as stdin
		if stdin == b {
			inBy := &bytes.Buffer{}
			for idx := l1; idx <= l2; idx++ {
				inBy.Write([]byte(b.getLine(idx)))
				inBy.Write([]byte("\n"))
			}
			stdin = inBy
		}

		shellCmd := strings.TrimSpace(command)

		// hide all escaped '%'
		shellCmd = strings.ReplaceAll(shellCmd, `\%`, string(rune(26)))
		// replace '%' with filename
		shellCmd = strings.ReplaceAll(shellCmd, `%`, b.getFilename())
		// put the '%' back
		shellCmd = strings.ReplaceAll(shellCmd, string(rune(26)), `%`)

		buf := bufio.NewScanner(strings.NewReader(shellCmd))
		buf.Split(bufio.ScanWords)

		var args []string
		for buf.Scan() {
			args = append(args, os.ExpandEnv(buf.Text()))
		}

		cmd := exec.Command(args[0], args[1:]...)

		var outBy bytes.Buffer

		// troubleshoot panics
		// var errBy bytes.Buffer
		// defer func() {
		// 	if err := recover(); err != nil {
		// 		stderr.Println(err, outBy.Len(), errBy.Len())
		// 	}
		// }()

		cmd.Stdin = stdin
		cmd.Stdout = &outBy
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		outBy.WriteByte('\n') // shim a newline because scripts that do not end with a newline cause a panic
		io.Copy(stdout, &outBy)

		return "!", err
	}
}

func doSetFilename(b buffer, fs FileSystem, filename string) (string, error) {
	if filename == "" {
		return b.getFilename(), nil
	}

	path, err := fs.Abs(filename)
	if err != nil {
		return "", err
	}

	b.setFilename(path)

	wd, err := os.Getwd()
	if err != nil {
		return path, nil
	}

	return filepath.Rel(wd, path)
	// dir := filepath.Dir(path)
	// if dir == "." {
	// 	path = filepath.Join(dir, path)
	// }

	// return path, err
}

func doSetMarkLine(b buffer, l1, l2 int, arg string) error {
	var (
		mk  rune
		err error
	)

	if arg == "" {
		mk = null
	} else {
		mk, err = firstRune(arg)
		if err != nil {
			return err
		}
	}

	for idx := l1; idx <= l2; idx++ {
		b.putMark(idx, mk)
	}
	return nil
}

// func doGetMarkedLines(out io.Writer, b buffer, m string) error {
// 	mk := null
// 	if len(m) > 0 {
// 		mk = rune(m[0])
// 	}

// 	return doMapMarkedLines(b, mk, func(b buffer, idx int) error {
// 		_, err := fmt.Fprintf(out, "%2d) %s", idx, b.getLine(idx))
// 		return err
// 	})

// }

func doGetNextMatchedLine(out io.Writer, b buffer, ser search) error {
	if len(ser.pattern) == 0 { // no pattern means to repeat the last search
		return errors.New("empty pattern")
	}

	re, err := regexp.Compile(ser.pattern)
	if err != nil {
		return err
	}

	scan := b.scanForward(b.nextLine(b.getCurline()), b.getLastline())
	if ser.reverse {
		scan = b.scanReverse(b.prevLine(b.getCurline()), b.getLastline())
	}

	for {
		i, ok := scan()
		if !ok {
			break
		}

		if re.MatchString(b.getLine(i)) {
			fmt.Fprintf(out, "%2d) %s", i, b.getLine(i))
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

func doDebug(b buffer, cache *cache) (string, error) {
	if allowDebug != `true` {
		return "", errors.New("debugging is not enabled; recompile to enable")
	}

	var rtn strings.Builder

	rtn.WriteString(getBuildInfo().String())
	rtn.WriteString(cache.String())
	rtn.WriteString(b.String())
	return rtn.String(), nil
}
