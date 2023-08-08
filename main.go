package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/henderjon/logger"
)

const prompt = ":"

var (
	stderr = logger.NewDropLogger(os.Stderr)
	// stdin   = bufio.NewScanner(os.Stdin)
	errQuit = errors.New("goodbye")
	pager   = 5
)

func main() {
	opts := getParams()
	var cursave = 1
	b := newMemoryBuf(opts.general.filename)
	b = fillDemo(b)

	cmdParser := &parser{}
	input := getInput(os.Stdin, os.Stdout)
	for { // main loop
		line, err := input(prompt)
		if err != nil {
			break
		}

		cursave = b.getCurline()

		cmd, err := cmdParser.run(line)
		if cmd == nil || err != nil {
			fmt.Fprintf(os.Stdout, "invalid command; %s\n", err.Error())
			b.setCurline(cursave)
			continue
		}

		if contains(string(prefixes), cmd.globalPrefix) {
			err = doGlob(*cmd, b, input)
			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
			}
		} else {
			err = doCmd(*cmd, b, input)
			switch true {
			case err == errQuit:
				fmt.Fprintln(os.Stdout, err.Error())
				os.Exit(0)
			case err != nil:
				fmt.Fprintln(os.Stdout, err.Error())
			}
		}
		// fmt.Fprint(os.Stdout, prompt)
	}
}

func doCmd(cmd command, b buffer, input interactor) error {
	var err error

	// some commands do not require addresses
	switch cmd.action {
	case quitAction:
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
		return doSimpleReplace(b, line1, line2, cmd.pattern, cmd.substitution, cmd.replaceNum)
	case regexReplaceAction:
		return doRegexReplace(b, line1, line2, cmd.pattern, cmd.substitution, cmd.replaceNum)
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
		return doToggleMarkLine(b, line1, line2)
	case getMarkAction:
		return doGetMarkedLine(b)
	case searchAction:
		return doGetNextMatchedLine(b, cmd.addrPattern, true)
	case searchRevAction:
		return doGetNextMatchedLine(b, cmd.addrPattern, false)
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

// doGlob is *big* because we're not using globals and it's called from a scope where it doesn't share information like the original implementation
func doGlob(cmd command, b buffer, input interactor) error {
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

	// needed later to restore cursor after glob
	cursave := b.getCurline()

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
			b.putMark(i, true)
			continue
		}
		b.putMark(i, false)
	}

	// scan will loop once for every line even if the action is destructive so it can lap itself
	// this shouldn't be an issue if we're handling getMark()s well and restoring curline when we're done
	scan = b.scanForward(line1, numLines)
	for {
		i, ok := scan()
		if !ok {
			break
		}

		if !b.getMark(i) || contains(string(prefixes), cmd.action) {
			continue
		}

		cmd.addrStart = ""
		cmd.addrEnd = ""
		b.putMark(i, false)
		b.setCurline(i)
		doCmd(cmd, b, input)
		b.setCurline(i)
	}

	b.setCurline(b.nextLine(cursave))
	return nil

	// loop over buffer, mark lines the match in order to keep track of what's been done because doCmd/do* can alter the buffer
	// loop over buffer, execute command on each marked line
}
