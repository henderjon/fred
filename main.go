package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/henderjon/logger"
)

const prompt = ":"

var (
	stderr  = logger.NewDropLogger(os.Stderr)
	stdin   = bufio.NewScanner(os.Stdin)
	errQuit = errors.New("goodbye")
	pager   = 0
)

func main() {
	b := newMemoryBuf("")
	b = fillDemo(b)

	cmdParser := &parser{}
	fmt.Fprint(os.Stdout, prompt)
	for stdin.Scan() { // main loop
		err := stdin.Err()
		if err != nil {
			break
		}

		cmd, err := cmdParser.run(stdin.Text())
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
		}

		if cmd == nil {
			fmt.Println("invalid command")
			fmt.Fprint(os.Stdout, prompt)
			continue
		}

		// cursave := b.curline
		if cmd.globalPrefix {
			// doCmd over range of lines
		} else {
			err := doCmd(*cmd, b)
			switch true {
			case err == errQuit:
				fmt.Fprintln(os.Stdout, err.Error())
				os.Exit(0)
			case err != nil:
				fmt.Fprintln(os.Stdout, err.Error())
			}
		}
		fmt.Fprint(os.Stdout, prompt)
	}
}

func doCmd(cmd command, b buffer) error {
	var err error

	// some commands do not require addresses
	switch cmd.action {
	case quitAction:
		return errQuit
	}

	// some commands require addresses
	// line1, line2, err := b.defaultLines(cmd.addrStart, cmd.addrEnd)
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
		return doAppend(b, line1)
	case insertAction:
		return doInsert(b, line1)
	case deleteAction:
		return doDelete(b, line1, line2)
	case changeAction:
		return doChange(b, line1, line2)
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
	}

	stderr.Log(line1, line2)
	stderr.Log(cmd)

	return err
}
