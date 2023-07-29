package main

import "errors"

const (
	eof   = -1
	unset = -1
	zero  = "1"
	end   = "$"
)

var (
	errQuit  = errors.New("goodbye")
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
	// firstLine = byte('0')
	// lastLine = byte('$')
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
	printAction        = rune('p')
	printNumsAction    = rune('n')
	printLiteralAction = rune('l')
	quitAction         = rune('q')
	appendAction       = rune('a')
	deleteAction       = rune('d')
	insertAction       = rune('i')
	changeAction       = rune('c')
	eqAction           = rune('=')
	moveAction         = rune('m')
	mirrorAction       = rune('M')
	copyAction         = rune('k')
	// markAction       = rune(' ')
	searchAction        = rune('/') // /re/... establishes the ADDRESSES for the lines against which to execute cmd
	gSearchAction       = rune('g') // g/re/p is the glob prefix which means we use the pattern to print every line that matches [gPREFIX]
	gReplaceAction      = rune('g') // s/foo/bar/g is the glob prefix which means we replace ALL the matches not just the first
	transliterateAction = rune('t')
	scrollAction        = rune('z')
	joinAction          = rune('j')
	simpleReplaceAction = rune('s') // Substitute CMD
	regexReplaceAction  = rune('S')
	editAction          = rune('e') // Edit command
	fileAction          = rune('f') // File command
	readAction          = rune('r') // read [file] command
	writeAction         = rune('w') // write [file] command
)

var cmds = []rune{
	printAction,
	printNumsAction,
	printLiteralAction,
	quitAction,
	appendAction,
	deleteAction,
	insertAction,
	changeAction,
	eqAction,
	moveAction,
	mirrorAction,
	copyAction,
	searchAction,
	transliterateAction,
	scrollAction,
	joinAction,
	simpleReplaceAction,
	regexReplaceAction,
	editAction,
	fileAction,
	readAction,
	writeAction,
}
