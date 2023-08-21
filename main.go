package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/henderjon/logger"
)

var (
	stderr         = logger.NewDropLogger(os.Stderr)
	errQuit        = errors.New("goodbye")
	errDirtyBuffer = errors.New("you have unsaved changes; use Q to quit without saving")
	pager          = 0
)

func main() {
	opts := getParams()
	pager = opts.general.pager
	b := newMemoryBuf()
	cache := &cache{}

	inout, destructor := newClassicTerm(os.Stdin, os.Stdout)
	if !opts.general.classic { // if we switch these, the terminal gets stuck in raw
		inout, destructor = newLocalTerm(os.Stdin, os.Stdout)
	}
	defer destructor()

	if len(opts.general.filename) > 0 {
		numbts, err := doReadFile(b, b.getCurline(), opts.general.filename)
		if err != nil {
			inout.println(err.Error())
			return
		}
		inout.println(numbts)
		b.setDirty(false) // loading the file on init isn't *actually* dirty
	}

	cmdParser := &parser{}
	for { // main loop
		var msg string
		line, err := inout.input(opts.general.prompt)
		if err != nil {
			break
		}

		cursave := b.getCurline()

		cmd, err := cmdParser.run(line)
		if cmd == nil || err != nil {
			inout.printf("invalid command; %s\n", err.Error())
			b.setCurline(cursave)
			continue
		}

		// regular 'g' or 'v'
		if contains(string(globsPre), cmd.globalPrefix) {
			err = doGlob(*cmd, b, inout, cache)
			if err != nil {
				inout.println(err.Error())
			}
			continue
		}

		// interactive 'G' or 'V'
		if contains(string(intrGlobsPre), cmd.globalPrefix) {
			err = doInteractiveGlob(*cmd, b, inout, cache, opts.general.prompt)
			if err != nil {
				inout.println(err.Error())
			}
			continue
		}

		// TODO: doCmd should return a string and an err and all our single line printing could be here which will avoid having to inject a printer

		msg, err = doCmd(*cmd, b, inout, cache)
		switch true {
		case cmd.subCommand == quitAction:
			if b.isDirty() {
				inout.println(errDirtyBuffer)
				continue
			}
			inout.println(errQuit)
			return
		case err == errQuit:
			inout.println(err.Error())
			return
		case err != nil:
			inout.println(err.Error())
		case msg != "":
			inout.println(msg)
		}
	}
}

func doCmd(cmd command, b buffer, inout termio, cache *cache) (string, error) {
	var err error

	// some commands do not require addresses
	switch cmd.action {
	case reallyQuitAction:
		b.setDirty(false)
		return "", errQuit
	case quitAction:
		if b.isDirty() {
			return "", errDirtyBuffer
		}
		return "", errQuit
	}

	// some commands require addresses
	line1, line2, err := b.defLines(cmd.addrStart, cmd.addrEnd, cmd.addrIncr, b.getCurline(), b.getCurline())
	if err != nil {
		return "", err
	}

	switch cmd.action {
	case helpAction:
		flag.Usage()
		return "", nil
	case 0: // rune default (empty action)
		return "", doPrint(inout, b, line1, line2, pager, printTypeNum)
	case eqAction:
		return doPrintAddress(b, line2)
	case printAction:
		return "", doPrint(inout, b, line1, line2, pager, printTypeReg)
	case printNumsAction:
		return "", doPrint(inout, b, line1, line2, pager, printTypeNum)
	case printLiteralAction:
		return "", doPrint(inout, b, line1, line2, pager, printTypeLit)
	case appendAction:
		return "", doAppend(inout, b, line1)
	case insertAction:
		return "", doInsert(inout, b, line1)
	case deleteAction:
		return "", doDelete(b, line1, line2)
	case changeAction:
		return "", doChange(inout, b, line1, line2)
	case moveAction:
		return "", doMove(b, line1, line2, cmd.destination)
	case copyAction:
		return "", doCopyNPaste(b, line1, line2, cmd.destination)
	case simpleReplaceAction:
		return "", doSimpleReplace(b, line1, line2, cmd.pattern, cmd.substitution, cmd.replaceNum, cache)
	case regexReplaceAction:
		return "", doRegexReplace(b, line1, line2, cmd.pattern, cmd.substitution, cmd.replaceNum, cache)
	case breakAction:
		return "", doBreakLines(b, line1, line2, cmd.pattern, cmd.substitution, cmd.replaceNum, cache)
	case joinAction:
		return "", doJoinLines(b, line1, line2, cmd.pattern)
	case transliterateAction:
		return "", doTransliterate(b, line1, line2, cmd.pattern, cmd.substitution)
	case mirrorAction:
		return "", doMirrorLines(b, line1, line2)
	case setPagerAction:
		return setPager(&pager, cmd.destination)
	case shellAction:
		return doExternalShell(b, line1, line2, cmd.argument)(false, os.Stdout)
	case filenameAction:
		return doSetFilename(b, cmd.argument)
	case putMarkAction:
		return "", doSetMarkLine(b, line1, line2, cmd.argument)
	case getMarkAction:
		return "", doGetMarkedLine(inout, b, cmd.argument)
	case searchAction:
		return "", doGetNextMatchedLine(inout, b, cmd.addrPattern, true, cache)
	case searchRevAction:
		return "", doGetNextMatchedLine(inout, b, cmd.addrPattern, false, cache)
	case reallyEditAction:
		b.setDirty(false)
		fallthrough // 'E' is exactly like edit but ignore the unsaved changes warning.
		// generally speaking "fallthrough" should be avoided, but these two commands are almost identical
	case editAction: // read into the current buffer either shell output or a file
		if err = clearBuffer(b); err != nil {
			return "", err
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
		return doWriteFile(inout, b, line1, line2, cmd.argument)
	}

	stderr.Log(line1, line2)
	stderr.Log(cmd)

	return "", err
}

// doGlob is *big* because we're not using globals and it's called from a scope where it doesn't share information like the original implementation
func doGlob(cmd command, b buffer, inout termio, cache *cache) error {
	var err error

	// some commands require addresses
	line1, line2, err := b.defLines(cmd.addrStart, cmd.addrEnd, cmd.addrIncr, b.getCurline(), b.getCurline())
	if err != nil {
		return err
	}

	// 'v' & 'V' do inverted search but are "global prefixes"
	invertSearch := contains(string([]rune{globalNegSearchAction, globalNegIntSearchAction}), cmd.globalPrefix)
	numLines := line2 - line1
	if numLines <= 0 { // I like big loops
		numLines = b.getNumLines()
	}

	err = doMarkLines(b, line1, numLines, cmd.addrPattern, invertSearch)
	if err != nil {
		return err
	}

	// needed later to restore cursor after glob
	cursave := b.getCurline()

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
		doCmd(cmd, b, inout, cache)
		b.setCurline(i)
	}

	b.setCurline(b.nextLine(cursave))
	return nil

	// loop over buffer, mark lines the match in order to keep track of what's been done because doCmd/do* can alter the buffer
	// loop over buffer, execute command on each marked line
}

func doInteractiveGlob(cmd command, b buffer, inout termio, cache *cache, prompt string) error {
	var err error

	// some commands require addresses
	line1, line2, err := b.defLines(cmd.addrStart, cmd.addrEnd, cmd.addrIncr, b.getCurline(), b.getCurline())
	if err != nil {
		return err
	}

	// 'v' & 'V' do inverted search but are "global prefixes"
	invertSearch := contains(string([]rune{globalNegSearchAction, globalNegIntSearchAction}), cmd.globalPrefix)
	numLines := line2 - line1
	if numLines <= 0 { // I like big loops
		numLines = b.getNumLines()
	}

	err = doMarkLines(b, line1, numLines, cmd.addrPattern, invertSearch)
	if err != nil {
		return err
	}

	// needed later to restore cursor after glob
	cursave := b.getCurline()

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

		inout.println(".. " + b.getLine(i))

		stop := false
		for !stop {
			line, err := inout.input(fmt.Sprintf(".. %s", prompt))
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

			_, err = doCmd(*cmd, b, inout, cache)
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
