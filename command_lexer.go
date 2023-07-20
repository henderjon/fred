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
	state stateFn   // the next lexing function to enter
	pos   int       // current position in the input.
	start int       // start position of this item.
	width int       // width of last rune read from input.
	items chan item // channel of scanned items.
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
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
	r := strings.ContainsRune(valid, l.next())
	// if !r {
	// 	l.backup()
	// }
	return r
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
	}
	l.backup() // we're always going to move one too many but accept() already backs up
	return i > 0
}

// lineNumber reports which line we're on. Doing it this way
// means we don't have to worry about peek double counting.
func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.pos], "\n")
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
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
	// l.emit(itemEOF)
	close(l.items)
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item, 2), // Two items sufficient.
	}
	return l
}
