package main

const (
	eof = -1
)

var (
	helpAction               = rune('h')
	printAction              = rune('p')
	printNumsAction          = rune('n')
	printLiteralAction       = rune('l')
	quitAction               = rune('q')
	reallyQuitAction         = rune('Q')
	appendAction             = rune('a')
	deleteAction             = rune('d')
	insertAction             = rune('i')
	changeAction             = rune('c')
	eqAction                 = rune('=')
	moveAction               = rune('m')
	mirrorAction             = rune('M')
	copyAction               = rune('k')
	putMarkAction            = rune('\'') // marks the current line
	getMarkAction            = rune('"')  // gets the next marked line
	searchAction             = rune('/')  // /re/... establishes the ADDRESSES for the lines against which to execute cmd moving forward
	searchRevAction          = rune('\\') // \re\... establishes the ADDRESSES for the lines against which to execute cmd moving backward
	globalSearchAction       = rune('g')  // g/re/p is the glob prefix which means we use the pattern to print every line that matches [gPREFIX]
	globalIntSearchAction    = rune('G')  // g/re/p is the glob prefix which means we use the pattern to print every line that matches [gPREFIX]
	globalNegSearchAction    = rune('v')  // v/re/p is the glob prefix which means we use the pattern to print every line that doesn't match [gPREFIX]
	globalNegIntSearchAction = rune('V')  // v/re/p is the glob prefix which means we use the pattern to print every line that doesn't match [gPREFIX]
	globalReplaceAction      = rune('g')  // s/foo/bar/g is the glob suffix which means we replace ALL the matches within a line not just the first
	transliterateAction      = rune('t')
	setPagerAction           = rune('z')
	joinAction               = rune('j')
	breakAction              = rune('b')
	simpleReplaceAction      = rune('s') // Substitute CMD
	regexReplaceAction       = rune('S')
	editAction               = rune('e') // Edit command
	reallyEditAction         = rune('E') // Edit command
	filenameAction           = rune('f') // File command
	readAction               = rune('r') // read [file] command
	writeAction              = rune('w') // write [file] command
	superWriteAction         = rune('W') // write [file] command
	shellAction              = rune('!')
)

var cmds = []rune{
	helpAction,
	printAction,
	printNumsAction,
	printLiteralAction,
	quitAction,
	reallyQuitAction,
	appendAction,
	deleteAction,
	insertAction,
	changeAction,
	eqAction,
	moveAction,
	mirrorAction,
	copyAction,
	putMarkAction,
	getMarkAction,
	searchAction,
	searchRevAction,
	transliterateAction,
	setPagerAction,
	joinAction,
	breakAction,
	simpleReplaceAction,
	regexReplaceAction,
	editAction,
	reallyEditAction,
	filenameAction,
	readAction,
	writeAction,
	superWriteAction,
	shellAction,
}

var globsPre = []rune{
	globalSearchAction,
	globalNegSearchAction,
}

var intrGlobsPre = []rune{
	globalIntSearchAction,
	globalNegIntSearchAction,
}
