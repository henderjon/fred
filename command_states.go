package main

import (
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
		switch true {
		// case r == eof:
		// 	return nil
		case isSpace(r):
			l.ignore()
		case r == '+' || r == '-' || ('0' <= r && r <= '9'):
			l.backup() // TODO: this might be a clue ...
			// stderr.Fatal(l.current())
			return lexAddress
		case r == gSearchAction:
			l.emit(itemGlobalPrefix)
			delim := l.next()
			l.ignore() // ignore the delim
			return lexPattern(delim, itemPattern)
		case r == searchAction:
			l.ignore() // ignore the delim
			return lexPattern(r, itemPattern)
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
	if !l.accept("+-") {
		l.backup()
	}

	digits := "0123456789"
	if l.acceptRun(digits) {
		l.emit(itemAddress)
	} else {
		return l.errorf("invalid or missing address: %s", l.current())
	}

	return lexDef
}

// lexCommand checks a run for being a valid command
func lexAction(l *lexer) stateFn {
	if l.acceptRun(string(cmds)) {
		// cmd := l.current()
		l.emit(itemAction)
		// // some actions need more information
		// switch cmd {
		// case string(moveAction):
		// 	lexDestination(l)
		// case string(copyAction):
		// 	lexDestination(l)
		// }
		return nil
	}

	return l.errorf("unknown command: %s", l.current())
}

// lexPattern checks for the regex pattern for 'g' and 's'
func lexPattern(delim rune, t itemType) stateFn {
	return stateFn(func(l *lexer) stateFn {
		// reject empty patterns

		if !l.acceptUntil(string(delim)) {
			return l.errorf("empty pattern or missing delim")
		}

		if delim == l.peek() {
			l.emit(t)
			l.acceptRun(string(delim)) // TODO: consuming it here ... how does that screw with our substitutions?
			l.ignore()
		} else {
			return l.errorf("missing the closing delim")
		}

		return lexDef
	})
}

// lexDestination checks for the trailing address for actions such as 'm' and 'k'
// func lexDestination(l *lexer) stateFn {
// 	// Optional leading sign.
// 	l.accept("+-")
// 	digits := "0123456789"
// 	if l.acceptRun(digits) {
// 		l.emit(itemDestination)
// 	} else {
// 		return l.errorf("current command requires a destination address")
// 	}

// 	return lexDef
// }
