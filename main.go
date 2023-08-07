package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/henderjon/logger"
)

const prompt = ":"

var (
	stderr = logger.NewDropLogger(os.Stderr)
	// stdin   = bufio.NewScanner(os.Stdin)
	errQuit = errors.New("goodbye")
	pager   = 0
)

func main() {
	var cursave = 1
	b := newMemoryBuf("")
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

		if cmd.globalPrefix {
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
		if cmd.subCommand == string(shellAction) {
			b.setCurline(line1)
			return doExternalShell(b, line1, line2, cmd.argument)(false, b)
		}
		return doReadFile(b, line1, cmd.argument)
	case writeAction: // write the current buffer to either shell (stdin) or a file
		if cmd.subCommand == string(shellAction) {
			b.setCurline(line1)
			return doExternalShell(b, line1, line2, cmd.argument)(true, os.Stdout)
		}
		return doWriteFile(b, line1, line2, cmd.argument)
	}

	stderr.Log(line1, line2)
	stderr.Log(cmd)

	return err
}

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

	scan := b.scanForward(line1, line2)
	for {
		i, ok := scan()
		if !ok {
			break
		}

		if re.MatchString(b.getLine(i)) {
			b.putMark(i, true)
			continue
		}
		b.putMark(i, false)
	}

	scan = b.scanForward(line1, line2)
	for {
		i, ok := scan()
		if !ok {
			break
		}

		if !b.getMark(i) || strings.ContainsRune(string([]rune{globalSearchAction}), cmd.action) {
			continue
		}

		cmd.addrStart = ""
		cmd.addrEnd = ""
		b.setCurline(i)
		doCmd(cmd, b, input)
		b.putMark(i, false)
		b.setCurline(i)
	}

	return nil

	//1,2g/[1-5]/s/or/%/g
	// loop over buffer, mark lines the match in order to keep track of what's been done because doCmd/do* can alter the buffer
	// loop over buffer, execute command on each marked line
}
