package main

import (
	"unicode"
)

// [num][range][num][action][delim][pattern][delim][pattern][delim][additional]

// lexDef is the default lex state
func lexDef(l *lexer) stateFn {
	for {
		if l.pos > l.start {
			l.emit(itemEOF)
			break
		}
		r := l.next()
		switch true {
		// case r == eof:
		// 	return nil
		case isSpace(r):
			l.ignore()
		case r == '+' || r == '-' || ('0' <= r && r <= '9'):
			l.backup()
			return lexAddress
		case isAlpha(r):
			l.backup()
			return lexAction
		case r == ',':
			l.emit(itemRange)
		case r == eof:
			l.emit(itemEOF)
			break
		default:
			stderr.Log(string(r))
			return lexErr
			// return nil //l.errorf("unrecognized character in action: %#U", r)
		}
	}
	return nil
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

// isAlphaNumeric reports whether r is alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isAlpha reports whether r is alphabetic
func isAlpha(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func lexErr(l *lexer) stateFn {
	l.emit(itemError)
	return nil
}

// lexAddress scans a number: decimal, octal, hex, float, or imaginary.  This
// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
// and "089" - but when it's wrong the input is invalid and the parser (via
// strconv) will notice.
func lexAddress(l *lexer) stateFn {
	// Optional leading sign.
	l.accept("+-")
	digits := "0123456789"
	if l.acceptRun(digits) {
		l.emit(itemNumber)
	} else {
		l.emit(itemEmpty)
	}

	return lexDef
}

// lexCommand checks a run for being a valid command
func lexAction(l *lexer) stateFn {
	if l.accept(string(cmds)) {
		cmd := l.current()
		l.emit(itemAction)
		// some actions need more information
		switch cmd {
		case string(moveAction):
			lexDestination(l)
		case string(copyAction):
			lexDestination(l)
		case string(searchAction):
			return lexPattern(l)
		case string(substituteAction):
			lexPattern(l)
			return lexSubstitution(l)
		}
		return nil
	} else {
		// if we got a letter but that letter isn't a command ...
		l.emit(itemUnknownCommand)
	}
	return lexDef
}

// lexPattern checks for the regex pattern for 'g' and 's'
func lexPattern(l *lexer) stateFn {
	delim := l.next()
	l.ignore()

	// reject empty patterns
	if !l.acceptUntil(string(delim)) {
		l.emit(itemEmptyPattern)
	}
	if delim == l.peek() {
		l.emit(itemPattern)
	} else {
		l.emit(itemMissingDelim)
	}

	return lexDef
}

// lexSubstitution checks for the replacement/substitution string for 's'
func lexSubstitution(l *lexer) stateFn {
	delim := l.next()
	l.ignore() // delim
	// we don't care if it's empty
	l.acceptUntil(string(delim))
	if delim == l.peek() {
		l.emit(itemSubstitution)
		l.accept(string(delim))
	} else {
		l.emit(itemMissingDelim)
	}

	return lexDef
}

// lexDestination checks for the trailing address for actions such as 'm' and 'k'
func lexDestination(l *lexer) stateFn {
	// Optional leading sign.
	l.accept("+-")
	digits := "0123456789"
	if l.acceptRun(digits) {
		l.emit(itemNumber)
	} else {
		l.emit(itemEmpty)
	}

	return lexDef
}

// lexFlag checks for the trailing 'g' suffix
func lexFlag(l *lexer) stateFn {
	if l.acceptRun("g") {
		l.emit(itemGlobalFlag)
	} else {
		l.emit(itemEmpty)
	}

	return lexDef
}
