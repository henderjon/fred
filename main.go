package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	stderr           = LogWriter{log.New(os.Stderr, "", 0)}
	errStop          = errors.New("stop")
	errQuit          = errors.New("goodbye")
	errDirtyBuffer   = errors.New("you have unsaved changes; use Q to quit without saving")
	errEmptyFilename = errors.New("empty filename")
	// errEd            = errors.New("the experienced user will know what is wrong")
)

func main() {
	opts := getParams()
	cache := &cache{}
	cache.setPager(opts.general.pager)
	b := newBuffer(osFS{})

	// create our shutdown listener and our terminal and load the file given via -file
	shd, inout := bootstrap(b, cache, opts)
	defer shd.Destructor()

	cmdParser := &parser{}
	for { // main loop
		var msg string
		line, err := inout.input(opts.general.prompt)
		if err != nil {
			fmt.Fprint(inout, err.Error())
			return
		}

		cursave := b.getCurline()

		cmd, err := cmdParser.run(line)
		if cmd == nil || err != nil {
			fmt.Fprintf(inout, "invalid command; %s", err.Error())
			b.setCurline(cursave)
			continue
		}

		if cmd.action == undoAction {
			if t, err := cache.unstageUndo(); err != nil {
				fmt.Fprint(inout, err.Error())
			} else {
				b = t.clone()
			}
			continue
		}

		msg, err = doCmd(*cmd, b, inout, osFS{}, cache) // NOTE: should doCmd return (string, error) or only (error)
		cache.stageUndo(b.clone())                      // cache confirms the incoming is different

		switch true {
		case cmd.subCommand == quitAction:
			if b.isDirty() {
				fmt.Fprint(inout, errDirtyBuffer)
				continue
			}
			if msg != "" { // if our command gave us a message, print it before quitting
				fmt.Fprint(inout, msg)
			}
			fmt.Fprint(inout, errQuit)
			return
		case err == errStop: // used by the interactive commands
			continue
		case err == errQuit:
			fmt.Fprint(inout, err.Error())
			return
		case err != nil:
			fmt.Fprint(inout, err.Error())
		case msg != "":
			fmt.Fprint(inout, msg)
		}
	}
}

func doCmd(cmd command, b buffer, inout termio, fsys FileSystem, cache *cache) (string, error) {
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

	// some commands require addresses; this takes the input and gives us sane lines
	line1, line2, err := b.defaultLines(cmd.addrStart, cmd.addrEnd, cmd.addrIncr, b.getCurline(), b.getCurline())
	if err != nil {
		return "", err
	}

	switch cmd.globalPrefix { // do bulk line actions
	case globalSearchAction:
		return "", doGlob(b, line1, line2, cmd, inout, fsys, cache)
	case globalNegSearchAction:
		return "", doGlob(b, line1, line2, cmd, inout, fsys, cache)
	case globalIntSearchAction:
		return "", doInteractiveGlob(b, line1, line2, cmd, inout, fsys, cache)
	case globalNegIntSearchAction:
		return "", doInteractiveGlob(b, line1, line2, cmd, inout, fsys, cache)
	case bulkMarkAction:
		return "", doManualBulk(b, cmd.addrPattern, cmd, inout, fsys, cache)
	}

	switch cmd.action {
	case helpAction:
		flag.Usage()
		return "", nil
	case 0: // rune default (empty action)
		return "", doPrint(inout, b, line1, line2, cache, cache.getPrintType())
	case eqAction:
		return doPrintAddress(b, line2)
	case printAction:
		return "", doPrint(inout, b, line1, line2, cache, printTypeReg)
	case printNumsAction:
		return "", doPrint(inout, b, line1, line2, cache, printTypeNum)
	case printLiteralAction:
		return "", doPrint(inout, b, line1, line2, cache, printTypeLit)
	case appendAction:
		return "", doAppend(inout, b, line1)
	case insertAction:
		return "", doInsert(inout, b, line1)
	case deleteAction:
		return "", doDelete(b, line1, line2)
	case changeAction:
		return "", doChange(inout, b, line1, line2)
	case historyAction:
		return "", inout.prtHistory(cmd.argument)
	case moveAction:
		return "", doMove(b, line1, line2, cmd.destination)
	case copyAction:
		return "", doCopyNPaste(b, line1, line2, cmd.destination)
	case simpleReplaceAction:
		return "", doSimpleReplace(b, line1, line2, cache.replace(cmd.pattern, cmd.substitution, cmd.replaceNum))
	case regexReplaceAction:
		return "", doRegexReplace(b, line1, line2, cache.replace(cmd.pattern, cmd.substitution, cmd.replaceNum))
	case breakAction:
		return "", doBreakLines(b, line1, line2, replace{cmd.pattern, cmd.substitution, cmd.replaceNum})
	case joinAction:
		return "", doJoinLines(b, line1, line2, cmd.pattern)
	case transliterateAction:
		return "", doTransliterate(b, line1, line2, cmd.pattern, cmd.substitution)
	case mirrorAction:
		return "", doMirrorLines(b, line1, line2)
	case setPagerAction:
		return doSetPager(cmd.destination, cache)
	case shellAction:
		return doExternalShell(b, line1, line2, cmd.argument)(nil, inout)
	case filenameAction:
		return doSetFilename(b, fsys, cmd.argument)
	case putMarkAction:
		return "", doSetMarkLine(b, line1, line2, cmd.argument)
	case searchAction:
		return "", doGetNextMatchedLine(inout, b, cache.search(cmd.addrPattern, false))
	case searchRevAction:
		return "", doGetNextMatchedLine(inout, b, cache.search(cmd.addrPattern, true))
	case reallyEditAction:
		b.setDirty(false)
		fallthrough // 'E' is exactly like edit but ignore the unsaved changes warning.
		// generally speaking "fallthrough" should be avoided, but these two commands are almost identical
	case editAction: // read into the current buffer either shell output or a file
		if b.isDirty() {
			return "", errDirtyBuffer
		}

		if b.getLastline() > 1 {
			if err = doDelete(b, 1, b.getLastline()); err != nil {
				return "", err
			}
		}

		if cmd.subCommand == shellAction {
			b.setCurline(line1)
			return doExternalShell(b, line1, line1, cmd.argument)(nil, b)
		}

		return doEditFile(b, fsys, cmd.argument)
	case readAction: // read into the current buffer either shell output or a file
		if cmd.subCommand == shellAction {
			b.setCurline(line1)
			return doExternalShell(b, line1, line2, cmd.argument)(nil, b)
		}

		return doReadFile(b, line1, fsys, cmd.argument)
	case writeAction: // write the current buffer to either shell (stdin) or a file
		if cmd.subCommand == shellAction {
			b.setCurline(line1)
			return doExternalShell(b, line1, line2, cmd.argument)(b, os.Stdout)
		}

		return doWriteFile(inout, b, line1, line2, fsys, cmd.argument)
	case debugAction:
		return doDebug(b, cache)
	}

	stderr.Println(line1, line2)
	stderr.Println(cmd)

	return "", err
}
