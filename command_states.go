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
			return lexNumber
		case isAlpha(r):
			l.backup()
			return lexAction
		case r == ',':
			l.emit(itemRange)
		default:
			// l.backup()
			return lexErr
			// return nil //l.errorf("unrecognized character in action: %#U", r)
		}
	}
	return nil
}

func lexErr(l *lexer) stateFn {
	l.emit(itemError)
	return nil
}

// lexNumber scans a number: decimal, octal, hex, float, or imaginary.  This
// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
// and "089" - but when it's wrong the input is invalid and the parser (via
// strconv) will notice.
func lexNumber(l *lexer) stateFn {
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

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r':
		return true
	}
	return false
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlpha(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

// lexCommand checks a run for being a valid command
func lexAction(l *lexer) stateFn {
	if l.accept(string(cmds)) {
		l.emit(itemAction)
	} else {
		// if we got a letter but that letter isn't a command ...
		l.emit(itemUnknownCommand) // TODO: do we need to support alpha delims?
	}
	return lexDef
}
