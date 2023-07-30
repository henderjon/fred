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
		// stderr.Log(string(r))
		switch true {
		case isSpace(r):
			l.ignore()
		case strings.ContainsRune("+-0123456789.$", r):
			l.backup()
			return lexAddress(itemAddress)(l)
		case r == shellAction:
			l.emit(itemAction)
			return lexArg
		case r == eqAction:
			l.emit(itemAction)
		case r == gSearchAction:
			l.emit(itemGlobalPrefix)
			delim := l.next()
			l.ignore() // ignore the delim
			return lexPattern(delim, itemAddressPattern)
		case r == searchAction:
			l.ignore() // ignore the delim
			return lexPattern(r, itemAddressPattern)
		case isAlpha(r):
			l.backup()
			return lexAction
		case r == ',':
			l.emit(itemRange)
		case r == eof:
			l.emit(itemEOF)
			return nil
		default:
			return lexErr
			// return nil //l.errorf("unrecognized character in action: %#U", r)
		}
	}
	// return nil
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

// isAlpha reports whether r is alphabetic
func isAlpha(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func lexErr(l *lexer) stateFn {
	l.emit(itemError)
	return nil
}

// lexAddress parses a value that represents an address in the command
func lexAddress(t itemType) stateFn {
	return stateFn(func(l *lexer) stateFn {
		l.acceptRun("+-") // TODO: will accept more than one ... :thinking_face:

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

	// consider arg vs pattern
	case l.acceptOne(string([]rune{joinAction})):
		l.emit(itemAction)
		delim := l.next()
		// stderr.Log(string(delim))
		l.ignore() // ignore the delim
		return lexPattern(delim, itemPattern)(l)
	case l.acceptOne(string([]rune{simpleReplaceAction, regexReplaceAction, transliterateAction})):
		l.emit(itemAction)
		delim := l.next()
		// stderr.Log(string(delim))
		l.ignore() // ignore the delim
		lexPattern(delim, itemPattern)(l)
		lexPattern(delim, itemSubstitution)(l)
		return lexReplaceNum(l)
	case l.acceptOne(string(cmds)):
		l.emit(itemAction)
		// some commands take a space and more info; later when I deviate from traditional ed, maybe take spaces all over
		if space := l.peek(); isSpace(space) || space == shellAction {
			return lexArg(l)
		}
		return lexDef
	}
	return l.errorf("unknown command: %s", l.current())
}

// lexPattern checks for the regex pattern for 'g' and 's'
func lexPattern(delim rune, t itemType) stateFn {
	var level int
	return stateFn(func(l *lexer) stateFn {
		// reject empty patterns
		if !l.acceptUntil(string(delim)) {
			return l.errorf("empty pattern or substitution or missing delim")
		}

		if delim == l.peek() {
			l.emit(t)
			l.acceptRun(string(delim))
			l.ignore()
			level++
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

func lexArg(l *lexer) stateFn {
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
