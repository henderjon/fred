package main

const (
	eof = -1
)

var (
	shellAction              = rune('!')
	getMarkAction            = rune('"')
	searchAction             = rune('/') // /re/... establishes the ADDRESSES for the lines against which to execute cmd moving forward
	putMarkAction            = rune('\'')
	searchRevAction          = rune('\\') // \re\... establishes the ADDRESSES for the lines against which to execute cmd moving backward
	eqAction                 = rune('=')
	appendAction             = rune('a')
	breakAction              = rune('b')
	changeAction             = rune('c')
	deleteAction             = rune('d')
	editAction               = rune('e') // Edit command
	reallyEditAction         = rune('E') // Edit command
	filenameAction           = rune('f')
	globalIntSearchAction    = rune('G')
	globalReplaceAction      = rune('g') // s/foo/bar/g is the glob suffix which means we replace ALL the matches within a line not just the first
	globalSearchAction       = rune('g') // g/re/p is the glob prefix which means we use the pattern to print every line that matches [gPREFIX]
	historyAction            = rune('H')
	helpAction               = rune('h')
	insertAction             = rune('i')
	joinAction               = rune('j')
	copyAction               = rune('k')
	printLiteralAction       = rune('l')
	mirrorAction             = rune('M')
	moveAction               = rune('m')
	printNumsAction          = rune('n')
	printAction              = rune('p')
	quitAction               = rune('q')
	reallyQuitAction         = rune('Q')
	readAction               = rune('r')
	simpleReplaceAction      = rune('s')
	regexReplaceAction       = rune('S')
	transliterateAction      = rune('t')
	undoAction               = rune('u')
	globalNegIntSearchAction = rune('V')
	globalNegSearchAction    = rune('v') // v/re/p is the glob prefix which means we use the pattern to print every line that doesn't match [gPREFIX]
	superWriteAction         = rune('W') // write [file] command
	writeAction              = rune('w') // write [file] command
	setPagerAction           = rune('z')
	setColumnAction          = rune('^')
	printColumnAction        = rune('|')
)

var cmds = []rune{
	shellAction,
	getMarkAction,
	searchAction,
	putMarkAction,
	searchRevAction,
	eqAction,
	appendAction,
	breakAction,
	changeAction,
	deleteAction,
	editAction,
	reallyEditAction,
	filenameAction,
	globalIntSearchAction,
	globalReplaceAction,
	globalSearchAction,
	historyAction,
	helpAction,
	insertAction,
	joinAction,
	copyAction,
	printLiteralAction,
	mirrorAction,
	moveAction,
	printNumsAction,
	printAction,
	quitAction,
	reallyQuitAction,
	readAction,
	simpleReplaceAction,
	regexReplaceAction,
	transliterateAction,
	undoAction,
	globalNegIntSearchAction,
	globalNegSearchAction,
	superWriteAction,
	writeAction,
	setPagerAction,
	setColumnAction,
	printColumnAction,
}

var globsPre = []rune{
	globalSearchAction,
	globalNegSearchAction,
}

var intrGlobsPre = []rune{
	globalIntSearchAction,
	globalNegIntSearchAction,
}
