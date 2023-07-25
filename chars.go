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
	null = byte('\x00')
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
	// gSUFFIX          = rune('g')      // s/../../g is the gflag that tells us to perform the substitution more than once on the line
	printAction     = rune('p') // Print CMD
	printNumsAction = rune('n') // Print CMD
	quitAction      = rune('q') // Quit CMD
	appendAction    = rune('a') // Append CMD
	deleteAction    = rune('d') // Delete CMD
	insertAction    = rune('i') // Insert CMD
	changeAction    = rune('c') // Change CMD
	eqAction        = rune('=') // Equal CMD
	moveAction      = rune('m') // Move CMD
	copyAction      = rune('k') // Copy CMD
	searchAction    = rune('/') // /re/... establishes the ADDRESSES for the lines against which to execute cmd
	gSearchAction   = rune('g') // g/re/p is the glob prefix which means we use the pattern to print every line that matches [gPREFIX]
	// TODO: there is a complexity around addressing lines via regex and then running a command ... there is the possibility of 3 patterns e.g. 10,20g/pattern/s/pattern/sub/
	substituteAction = rune('s') // Substitute CMD
	editAction       = rune('e') // Edit command
	fileAction       = rune('f') // File command
	readAction       = rune('r') // read [file] command
	writeAction      = rune('w') // write [file] command
)

var cmds = []rune{
	printAction,
	printNumsAction,
	quitAction,
	appendAction,
	deleteAction,
	insertAction,
	changeAction,
	eqAction,
	moveAction,
	copyAction,
	searchAction,
	substituteAction,
	editAction,
	fileAction,
	readAction,
	writeAction,
}

var destinationCmds = []rune{
	printAction,
	printNumsAction,
	moveAction,
	copyAction,
}
