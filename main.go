package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/henderjon/logger"
)

var (
	stderr = logger.NewDropLogger(os.Stderr)
	// stdin   = bufio.NewScanner(os.Stdin)
	errQuit        = errors.New("goodbye")
	errDirtyBuffer = errors.New("you have unsaved changes; use Q to quit without saving")
	pager          = 0
)

func main() {
	opts := getParams()
	pager = opts.general.pager
	b := newMemoryBuf()
	cache := &cache{}

	if len(opts.general.filename) > 0 {
		err := doReadFile(b, b.getCurline(), opts.general.filename)
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
			os.Exit(1)
		}
		b.setDirty(false) // loading the file on init isn't *actually* dirty
	}

	cmdParser := &parser{}
	input := getInput(os.Stdin, os.Stdout)
	for { // main loop
		line, err := input(opts.general.prompt)
		if err != nil {
			break
		}

		cursave := b.getCurline()

		cmd, err := cmdParser.run(line)
		if cmd == nil || err != nil {
			fmt.Fprintf(os.Stdout, "invalid command; %s\n", err.Error())
			b.setCurline(cursave)
			continue
		}

		// regular 'g' or 'v'
		if contains(string(globsPre), cmd.globalPrefix) {
			err = doGlobMarks(*cmd, b)
			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
				continue
			}

			err = doGlob(*cmd, b, input, cache)
			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
			}
			continue
		}

		// interactive 'G' or 'V'
		if contains(string(intrGlobsPre), cmd.globalPrefix) {
			err = doGlobMarks(*cmd, b)
			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
				continue
			}

			err = doInteractiveGlob(*cmd, b, input, cache, opts.general.prompt)
			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
			}
			continue
		}

		err = doCmd(*cmd, b, input, cache)
		switch true {
		case cmd.subCommand == quitAction:
			if b.isDirty() {
				fmt.Fprintln(os.Stdout, errDirtyBuffer)
				continue
			}
			fmt.Fprintln(os.Stdout, errQuit)
			os.Exit(0)
		case err == errQuit:
			fmt.Fprintln(os.Stdout, err.Error())
			os.Exit(0)
		case err != nil:
			fmt.Fprintln(os.Stdout, err.Error())
		}

	}
}

func doCmd(cmd command, b buffer, input interactor, cache *cache) error {
	var err error

	// some commands do not require addresses
	switch cmd.action {
	case reallyQuitAction:
		b.setDirty(false)
		return errQuit
	case quitAction:
		if b.isDirty() {
			return errDirtyBuffer
		}
		return errQuit
	}

	// some commands require addresses
	line1, line2, err := b.defLines(cmd.addrStart, cmd.addrEnd, b.getCurline(), b.getCurline())
	if err != nil {
		return err
	}

	switch cmd.action {
	case helpAction:
		flag.Usage()
		return nil
	case 0: // rune default (empty action)
		return doPrint(b, line1, line2, pager, printTypeNum)
	case eqAction:
		return doPrintAddress(b, line2)
	case printAction:
		return doPrint(b, line1, line2, pager, printTypeReg)
	case printNumsAction:
		return doPrint(b, line1, line2, pager, printTypeNum)
	case printLiteralAction:
		return doPrint(b, line1, line2, pager, printTypeLit)
	case appendAction:
		return doAppend(input, b, line1)
	case insertAction:
		return doInsert(input, b, line1)
	case deleteAction:
		return doDelete(b, line1, line2)
	case changeAction:
		return doChange(input, b, line1, line2)
	case moveAction:
		return doMove(b, line1, line2, cmd.destination)
	case copyAction:
		return doCopyNPaste(b, line1, line2, cmd.destination)
	case simpleReplaceAction:
		return doSimpleReplace(b, line1, line2, cmd.pattern, cmd.substitution, cmd.replaceNum, cache)
	case regexReplaceAction:
		return doRegexReplace(b, line1, line2, cmd.pattern, cmd.substitution, cmd.replaceNum, cache)
	case joinAction:
		return doJoinLines(b, line1, line2, cmd.pattern)
	case transliterateAction:
		return doTransliterate(b, line1, line2, cmd.pattern, cmd.substitution)
	case mirrorAction:
		return doMirrorLines(b, line1, line2)
	case setPagerAction:
		return setPager(&pager, cmd.destination)
	case shellAction:
		return doExternalShell(b, line1, line2, cmd.argument)(false, os.Stdout)
	case filenameAction:
		return doSetFilename(b, cmd.argument)
	case putMarkAction:
		return doSetMarkLine(b, line1, line2, cmd.argument)
	case getMarkAction:
		return doGetMarkedLine(b, cmd.argument)
	case searchAction:
		return doGetNextMatchedLine(b, cmd.addrPattern, true, cache)
	case searchRevAction:
		return doGetNextMatchedLine(b, cmd.addrPattern, false, cache)
	case editAction: // read into the current buffer either shell output or a file
		if err = clearBuffer(b); err != nil {
			return err
		}

		if cmd.subCommand == shellAction {
			b.setCurline(line1)
			return doExternalShell(b, line1, line1, cmd.argument)(false, b)
		}
		return doReadFile(b, line1, cmd.argument)
	case readAction: // read into the current buffer either shell output or a file
		if cmd.subCommand == shellAction {
			b.setCurline(line1)
			return doExternalShell(b, line1, line2, cmd.argument)(false, b)
		}
		return doReadFile(b, line1, cmd.argument)
	case writeAction: // write the current buffer to either shell (stdin) or a file
		if cmd.subCommand == shellAction {
			b.setCurline(line1)
			return doExternalShell(b, line1, line2, cmd.argument)(true, os.Stdout)
		}
		return doWriteFile(b, line1, line2, cmd.argument)
	}

	stderr.Log(line1, line2)
	stderr.Log(cmd)

	return err
}

func doGlobMarks(cmd command, b buffer) error {
	if len(cmd.addrPattern) == 0 {
		return fmt.Errorf("missing address pattern")
	}

	re, err := regexp.Compile(cmd.addrPattern)
	if err != nil {
		return err
	}

	// some commands require addresses
	line1, line2, err := b.defLines(cmd.addrStart, cmd.addrEnd, b.getCurline(), b.getCurline())
	if err != nil {
		return err
	}

	// 'v' & 'V' do inverted search but are "global prefixes"
	invertSearch := contains(string([]rune{globalNegSearchAction, globalNegIntSearchAction}), cmd.globalPrefix)

	// our scan takes an upper bound number of iterations
	numLines := line2 - line1
	if numLines <= 0 {
		numLines = b.getLastline() // all lines
	}

	scan := b.scanForward(line1, numLines)
	for {
		i, ok := scan()
		if !ok {
			break
		}

		if re.MatchString(b.getLine(i)) != invertSearch {
			b.putMark(i, mark)
			continue
		}
		// previously, we blanked every line's mark creating un/marked lines. If marks can be any rune, we only need to assert the mark of the lines we care about, right?
		// b.putMark(i, null)
	}
	return nil
}

// doGlob is *big* because we're not using globals and it's called from a scope where it doesn't share information like the original implementation
func doGlob(cmd command, b buffer, input interactor, cache *cache) error {
	// some commands require addresses
	line1, line2, err := b.defLines(cmd.addrStart, cmd.addrEnd, b.getCurline(), b.getCurline())
	if err != nil {
		return err
	}

	// needed later to restore cursor after glob
	cursave := b.getCurline()

	// our scan takes an upper bound number of iterations
	numLines := line2 - line1
	if numLines <= 0 {
		numLines = b.getLastline() // all lines
	}

	// scan will loop once for every line even if the action is destructive so it can lap itself
	// this shouldn't be an issue if we're handling getMark()s well and restoring curline when we're done
	scan := b.scanForward(line1, numLines)
	for {
		i, ok := scan()
		if !ok {
			break
		}

		if !b.hasMark(i, mark) ||
			contains(string(globsPre), cmd.action) ||
			contains(string(intrGlobsPre), cmd.action) {
			continue
		}

		cmd.addrStart = ""
		cmd.addrEnd = ""
		b.putMark(i, null)
		b.setCurline(i)
		doCmd(cmd, b, input, cache)
		b.setCurline(i)
	}

	b.setCurline(b.nextLine(cursave))
	return nil

	// loop over buffer, mark lines the match in order to keep track of what's been done because doCmd/do* can alter the buffer
	// loop over buffer, execute command on each marked line
}

func doInteractiveGlob(cmd command, b buffer, input interactor, cache *cache, prompt string) error {
	// some commands require addresses
	line1, line2, err := b.defLines(cmd.addrStart, cmd.addrEnd, b.getCurline(), b.getCurline())
	if err != nil {
		return err
	}

	// needed later to restore cursor after glob
	cursave := b.getCurline()

	// our scan takes an upper bound number of iterations
	numLines := line2 - line1
	if numLines <= 0 {
		numLines = b.getLastline() // all lines
	}

	// scan will loop once for every line even if the action is destructive so it can lap itself
	// this shouldn't be an issue if we're handling getMark()s well and restoring curline when we're done
	scan := b.scanForward(line1, numLines)
	for {
		i, ok := scan()
		if !ok {
			break
		}

		if !b.hasMark(i, mark) ||
			contains(string(globsPre), cmd.action) ||
			contains(string(intrGlobsPre), cmd.action) {
			continue
		}

		cmd.addrStart = ""
		cmd.addrEnd = ""
		b.putMark(i, null)
		b.setCurline(i)

		fmt.Fprintln(os.Stdout, ".. "+b.getLine(i))

		stop := false
		for !stop {
			line, err := input(fmt.Sprintf(".. %s", prompt))
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
				stop = true
				continue
			case cmd.action == reallyQuitAction:
				return nil
			}

			err = doCmd(*cmd, b, input, cache)
			switch true {
			case err != nil:
				stop = true
			case err == errQuit:
				fmt.Println(err)
			}
			b.setCurline(i)
		}
		b.setCurline(i)
	}

	b.setCurline(b.nextLine(cursave))
	return nil

	// loop over buffer, mark lines the match in order to keep track of what's been done because doCmd/do* can alter the buffer
	// loop over buffer, execute command on each marked line
}
