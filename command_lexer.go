package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name  string    // the name of the input; used only for error reports.
	input string    // the string being scanned.
	pos   int       // current position in the input.
	start int       // start position of this item.
	width int       // width of last rune read from input.
	items chan item // channel of scanned items.
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if len(l.input[l.pos:]) <= 0 {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// current shows the value to be emitted
func (l *lexer) current() string {
	return l.input[l.start:l.pos]
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	// l.next consumes a character
	return strings.ContainsRune(valid, l.next())
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptOne(valid string) bool {
	var i int
	if l.accept(valid) {
		i++
	} else {
		l.backup()
	}
	return i > 0
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) bool {
	var i int
	for ; l.accept(valid); i++ {
	}
	l.backup() // we're always going to move one too many but accept() already backs up
	return i > 0
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptUntil(invalid string) bool {
	var i int
	for ; !l.accept(invalid); i++ {
		if l.width == 0 {
			return false // if we run out of chars before we find our end
		}
	}
	l.backup() // we're always going to move one too many but accept() already backs up
	return i > 0
}

func (l *lexer) bleed() {
	var (
		i int
		n rune
	)
	for ; l.pos <= len(l.input); i++ {
		n = l.next()
		if n < 0 {
			break
		}
	}
	// l.backup() // we're always going to move one too many but accept() already backs up
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.run.
func (l *lexer) errorf(format string, args ...any) stateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
func (l *lexer) nextItem() item {
	return <-l.items
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for state := lexDef; state != nil; {
		state = state(l)
	}
	l.emit(itemEOF)
	close(l.items)
}

// lex creates a new scanner for the input string.
func lex(input, name string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item, 2), // Two items sufficient.
	}
	return l
}
