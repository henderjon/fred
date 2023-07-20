package main

import "errors"

const (
	eof   = -1
	unset = -1
	zero  = "0"
	end   = "-1"
)

var (
	errEOD   = errors.New("end of data")
	errBREAK = errors.New("break operation; no error")
	// line1        int // the first given line
	// line2        int // the second given line
	// nlines       int // the number of lines given
	// curline      int // the current line
	// lastline     int // last line of the buffer
	NULLBYTE = byte('\x00')
	// ENDOFSTRING  = NULLBYTE
	// DITTO        = byte('&')
	// CURRENT_LINE = byte('.')
	// LAST_LINE    = byte('$')
	// SCAN         = byte('/')
	// BACK_SCAN    = byte('\\')
	// ESCAPE       = byte('@')
	// NEWLINE      = byte('\n')
	// TAB          = byte('\t')
	// DASH         = byte('-')
	// BLANK        = byte(' ')
	// CLOSIZE      = 1
	// LITCHAR      = byte('c')
	// CLOSURE      = byte('*')
	// BOL          = byte('%')
	// EOL          = byte('$')
	// ANY          = byte('?')
	// CCL          = byte('[') // CCL == character class
	// CCLEnd       = byte(']')
	// NEGATE       = byte('^')
	// NCCL         = byte('!')
	// PERIOD       = byte('.') // Append CMD
	gPREFIX          = rune('g')      // g/../p is the glob prefix which means we use the pattern to print every line that matches
	gSUFFIX          = rune('g')      // s/../../g is the gflag that tells us to perform the substitution more than once on the line
	noAction         = rune(NULLBYTE) // NO ACTION
	printAction      = rune('p')      // Print CMD
	quitAction       = rune('q')      // Quit CMD
	appendAction     = rune('a')      // Append CMD
	deleteAction     = rune('d')      // Delete CMD
	insertAction     = rune('i')      // Insert CMD
	changeAction     = rune('c')      // Change CMD
	eqAction         = rune('=')      // Equal CMD
	moveAction       = rune('m')      // Move CMD
	copyAction       = rune('k')      // Copy CMD
	substituteAction = rune('s')      // Substitute CMD
	editAction       = rune('e')      // Edit command
	fileAction       = rune('f')      // File command
	readAction       = rune('r')      // read [file] command
	writeAction      = rune('w')      // write [file] command
)

var cmds = []rune{
	printAction,
	quitAction,
	appendAction,
	deleteAction,
	insertAction,
	changeAction,
	eqAction,
	moveAction,
	copyAction,
	substituteAction,
	editAction,
	fileAction,
	readAction,
	writeAction,
}
