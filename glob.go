// glob.go is a collection of the bulk line actions. The idea is that one
// function loops over all the lines and marks the relevant ones. A second loop
// then acts on each line that is marked.
package main

import (
	"fmt"
	"regexp"
)

type mapLine func(b buffer, idx int) error

// traverse the buffer and execute the given mapLine on each line with the given mark
func doMapMarkedLines(b buffer, m rune, fn mapLine) error {
	scan := b.scanForward(b.nextLine(b.getCurline()), b.getLastline())
	for {
		i, ok := scan()
		if !ok {
			break
		}

		if b.getMark(i) == m {
			b.putMark(i, null)

			err := fn(b, i)
			if err != nil {
				return err // short circuit the loop
			}

			b.setCurline(i)
		}
	}

	return nil
}

// doMarkLinesRegex walks the buffer and marks the lines matching `pattern` for further processing
func doMarkLinesRegex(b buffer, line1, line2 int, pattern string, invert bool) error {
	if len(pattern) == 0 {
		return fmt.Errorf("missing search pattern")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	numLines := line2 - line1
	if numLines <= 0 { // I like big loops
		numLines = b.getLastline()
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

// doGlob marks eash line according to `pattern` and executes `cmd` on that line
func doGlob(b buffer, line1, line2 int, cmd command, inout termio, fsys fileSystem, cache *cache) error {
	var err error

	// 'v' & 'V' do inverted search but are "global prefixes"
	err = doMarkLinesRegex(b, line1, line2, cmd.addrPattern, invertDirection(cmd.globalPrefix))
	if err != nil {
		return err
	}

	// needed later to restore cursor after glob
	cursave := b.getCurline()

	cmd.globalPrefix = null // blank the prefix so we're not recursing infinitely
	err = doMapMarkedLines(b, mark, globMapper(cmd, inout, fsys, cache))

	b.setCurline(b.nextLine(cursave))
	return err
}

// globMapper returns mapLine func that executes `cmd` on a line
func globMapper(cmd command, inout termio, fsys fileSystem, cache *cache) mapLine {
	return func(b buffer, idx int) error {
		if excludeFromGlob(cmd.action) {
			return nil
		}

		b.setCurline(idx)
		_, err := doCmd(cmd, b, inout, fsys, cache)
		return err
	}
}

// doGlob marks eash line according to `pattern` and prompts the user for commands
func doInteractiveGlob(b buffer, line1, line2 int, cmd command, inout termio, fsys fileSystem, cache *cache) error {
	var err error

	// 'v' & 'V' do inverted search but are "global prefixes"
	err = doMarkLinesRegex(b, line1, line2, cmd.addrPattern, invertDirection(cmd.globalPrefix))
	if err != nil {
		return err
	}

	// needed later to restore cursor after glob
	cursave := b.getCurline()

	cmd.globalPrefix = null // blank the prefix so we're not recursing infinitely
	err = doMapMarkedLines(b, mark, interactiveGlobMapper(cmd, inout, fsys, cache, fmt.Sprintf(".. %s", inout.getPrompt())))

	b.setCurline(b.nextLine(cursave))
	return err
}

// interactiveGlobMapper returns mapLine func that a user provided command on a line
func interactiveGlobMapper(cmd command, inout termio, fsys fileSystem, cache *cache, prompt string) mapLine {

	return func(b buffer, idx int) error {
		var err error
		if excludeFromGlob(cmd.action) {
			return nil
		}

		b.setCurline(idx)

		fmt.Fprint(inout, ".. "+b.getLine(idx))

		stop := false
		for !stop {
			line, err := inout.input(prompt)
			if err != nil {
				return err
			}

			cmd, err := (&parser{}).run(line)
			if err != nil {
				return err
			}

			switch true {
			case cmd.action == null: // the normal doCmd prints a line with no action; skip that behavior here
				continue
			case cmd.action == quitAction:
				stop = true // move to the next line
				continue
			case cmd.action == reallyQuitAction:
				return errStop
			}

			_, err = doCmd(*cmd, b, inout, fsys, cache)
			switch true {
			case err != nil:
				stop = true
			case err == errQuit:
				fmt.Println(err)
			}
			b.setCurline(idx)
		}
		return err
	}
}

// loop over buffer, mark lines the match in order to keep track of what's been done because doCmd/do* can alter the buffer
// loop over buffer, execute command on each marked line
