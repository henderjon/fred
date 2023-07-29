package main

const (
	eof = -1
)

var (
	printAction         = rune('p')
	printNumsAction     = rune('n')
	printLiteralAction  = rune('l')
	quitAction          = rune('q')
	appendAction        = rune('a')
	deleteAction        = rune('d')
	insertAction        = rune('i')
	changeAction        = rune('c')
	eqAction            = rune('=')
	moveAction          = rune('m')
	mirrorAction        = rune('M')
	copyAction          = rune('k')
	markAction          = rune('\'')
	searchAction        = rune('/') // /re/... establishes the ADDRESSES for the lines against which to execute cmd
	gSearchAction       = rune('g') // g/re/p is the glob prefix which means we use the pattern to print every line that matches [gPREFIX]
	gReplaceAction      = rune('g') // s/foo/bar/g is the glob suffix which means we replace ALL the matches within a line not just the first
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
	markAction,
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
