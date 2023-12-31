package main

import (
	"strings"
	"unicode"
)

// [num][range][num][action][delim][pattern][delim][pattern][delim][additional]

// lexDef is the default lex state
func lexDef(l *lexer) stateFn {
	for {
		// if l.pos > l.start {
		// 	l.emit(itemEOF)
		// 	break
		// }
		r := l.next()
		// stderr.Println(string(r))
		switch true {
		case isSpace(r):
			l.ignore()
		case strings.ContainsRune("+-0123456789.$", r):
			l.backup()
			return lexAddress(itemAddress)(l)
		case isIncr(r):
			l.emit(itemIncr)
			return lexAddress(itemAddress)(l)
		case r == shellAction:
			l.emit(itemAction)
			return lexArgs
		case r == historyAction:
			l.emit(itemAction)
			return lexArgs
		case r == eqAction:
			l.emit(itemAction)
		case isGlobalPrefix(r):
			l.emit(itemGlobalPrefix)
			delim := l.next()
			l.ignore() // ignore the delim
			return lexPattern(delim, itemAddressPattern)
		case isAction(r):
			l.backup()
			return lexAction
		case r == ',':
			l.emit(itemRange)
		case r == eof:
			l.emit(itemEOF)
			return nil
		default:
			l.errorf("unrecognized character in action: %s", string(r))
			// return lexErr
			return nil
		}
	}
	// return nil
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func isAction(r rune) bool {
	return strings.ContainsRune(allCmds(), r)
}

func isIncr(r rune) bool {
	return strings.ContainsRune("<>", r)
}

func isGlobalPrefix(r rune) bool {
	return strings.ContainsRune(string([]rune{
		globalSearchAction,
		globalIntSearchAction,
		globalNegSearchAction,
		globalNegIntSearchAction,
	}), r)
}

// lexAddress parses a value that represents an address in the command
func lexAddress(t itemType) stateFn {
	return stateFn(func(l *lexer) stateFn {
		l.acceptRun("+ -") // TODO: will accept more than one ... :thinking_face:

		switch true {
		default:
			l.emit(t)
			// return l.errorf("invalid or missing address/destination: %s", l.current())
		case l.acceptRun("$."):
			l.emit(t)
		case l.acceptRun("01234567890"):
			l.emit(t)
		}

		return lexDef
	})
}

// lexCommand checks a run for being a valid command
func lexAction(l *lexer) stateFn {
	switch true {
	// these commands need a destination
	case l.acceptOne(string([]rune{moveAction, copyAction, setPagerAction})):
		l.emit(itemAction)
		return lexAddress(itemDestination)(l)
	case l.acceptOne(string([]rune{searchAction})):
		l.emit(itemAction)
		return lexPattern(searchAction, itemAddressPattern)
	case l.acceptOne(string([]rune{searchRevAction})):
		l.emit(itemAction)
		return lexPattern(searchRevAction, itemAddressPattern)
	case l.acceptOne(string([]rune{joinAction, breakAction})): // TODO: join doesn't need a replace num... consider arg vs pattern ... ?
		l.emit(itemAction)
		delim := l.next()
		// stderr.Println(string(delim))
		l.ignore() // ignore the delim
		lexPattern(delim, itemPattern)(l)
		// lexPattern(delim, itemSubstitution)(l)
		return lexReplaceNum(l)
	case l.acceptOne(string([]rune{bulkMarkAction})):
		l.emit(itemGlobalPrefix)
		l.next() //  only take the next rune for an address pattern
		l.emit(itemAddressPattern)
		return lexDef
	case l.acceptOne(string([]rune{putMarkAction})):
		l.emit(itemAction)
		return lexArgs(l)
	case l.acceptOne(string([]rune{simpleReplaceAction, regexReplaceAction, transliterateAction})):
		l.emit(itemAction)
		delim := l.next()
		// stderr.Println(string(delim))
		l.ignore() // ignore the delim
		lexPattern(delim, itemPattern)(l)
		lexPattern(delim, itemSubstitution)(l)
		return lexReplaceNum(l)
	case l.acceptOne(allCmds()):
		l.emit(itemAction)
		// some commands take a space and more info; later when I deviate from
		// traditional ed, maybe take spaces all over
		if space := l.peek(); isSpace(space) || space == shellAction {
			return lexArgs(l)
		}
		return lexDef
	}
	return l.errorf("unknown command: %s", l.current())
}

// lexPattern checks for the regex pattern for 'g' and 's'
func lexPattern(delim rune, t itemType) stateFn {
	return stateFn(func(l *lexer) stateFn {
		// consume anything until next delim; allow empty patterns

		l.acceptUntil(string(delim))

		if delim == l.peek() {
			l.emit(t)
			l.acceptOne(string(delim))
			l.ignore()
		} else {
			return l.errorf("missing the closing delim")
		}

		return lexDef
	})
}

func lexReplaceNum(l *lexer) stateFn {
	if l.acceptRun("g") {
		l.emit(itemReplaceNum)
		return lexDef
	}

	digits := "0123456789"
	if l.acceptRun(digits) {
		l.emit(itemReplaceNum)
		return lexDef
	}

	return lexDef
}

func lexArgs(l *lexer) stateFn {
	if l.acceptRun(" ") {
		l.ignore()
	}

	if l.acceptRun(string(shellAction)) {
		l.emit(itemAction)
	}

	l.bleed()
	l.emit(itemArg)
	return lexDef
}

// func lexArg(l *lexer) stateFn {
// 	if l.acceptRun(" ") {
// 		l.ignore()
// 	}

// 	l.next() //  only take the next rune for an arg
// 	l.emit(itemArg)
// 	return lexDef
// }
