package main

import (
	"fmt"
	"io"
	"strings"
)

const (
	eof = -1
)

var (
	shellAction              = rune('!')
	bulkMarkAction           = rune('"')
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
	// setColumnAction          = rune('^')
	// printColumnAction        = rune('|')
	debugAction = rune('*')
)

type quickdoc struct {
	example     string
	explanation string
}

var cmds = map[rune]quickdoc{
	shellAction:              {`!exe`, "[shell] run $exe and print the output"},
	bulkMarkAction:           {`(,)"a cmd`, "[mark] print each line marked with the single character $a or do $cmd, if provided, for each line"},
	searchAction:             {`/regex/`, "[search] find the next line matching $regex"},
	putMarkAction:            {`(,)'a`, "[mark] mark the given lines with the single character $a"},
	searchRevAction:          {`\regex\`, "[search] find the previous line matching $regex"},
	eqAction:                 {`=`, "[print] print the address of the current line"},
	appendAction:             {`(.)a`, "[append] adds text to the buffer after the current line"},
	breakAction:              {`(,)b/regex/`, "[break] breaks the given lines into many lines at $regex"},
	changeAction:             {`(,)c`, "[change] replaces the given lines with user input"},
	deleteAction:             {`(,)d`, "[delete] removes the given lines from the buffer"},
	editAction:               {"e filename|!exe", "[edit] replaces the current buffer with $filename or the output of $exe"},
	filenameAction:           {"f name", "[file] set the current buffer to write to file $name"},
	reallyEditAction:         {"E file", "[edit] replaces the current buffer with $file without prompting to save changes"},
	globalIntSearchAction:    {"G/regex/", "[global] interactively take commands to run against all lines matching $regex"},
	globalReplaceAction:      {"", ""},
	globalSearchAction:       {"g/regex/cmd", "[global] run $cmd against each line matching $regex"},
	historyAction:            {"H n", "[history] print the last $n commands; disabled when compiled with '-tag readline'"},
	helpAction:               {"h", "[help] print the manual"},
	insertAction:             {"(.)i", "[insert] adds text to the buffer before the current line"},
	joinAction:               {"(,)j/sep/", "[join] join the given lines using $sep"},
	copyAction:               {"(,)k(.)", "[kopy/paste] copy the given lines to after the given destination"},
	printLiteralAction:       {"(,)l", "[literal] print the given lines but show special characters"},
	mirrorAction:             {"(,)M", "[mirror] reverse the order of the given lines"},
	moveAction:               {"(,)m(.)", "[move] move the given lines to after the given destination"},
	printNumsAction:          {"(,)n", "[number] print the given lines with line numbers in the left margin"},
	printAction:              {"(,)p", "[print]  print the given lines"},
	quitAction:               {"q", "[quit] exit the application"},
	reallyQuitAction:         {"Q", "[quit] exit the application without prompting to save changes"},
	readAction:               {"(.)r file|!exe", "[read] read the contents of $file or the output of $exe into the buffer after the given line"},
	simpleReplaceAction:      {"(,)s/pat/sub/n", "[substitute] replace $pat with $sub"},
	regexReplaceAction:       {"(,)S/regex/sub/n", "[substitute] replace $regex with $sub"},
	transliterateAction:      {"(,)t/find/repl/", "[transliterate] replace each character in $find with the corresponding character in $repl"},
	undoAction:               {"u", "[undo/redo] undo the last action; undo is it's own inverse"},
	globalNegIntSearchAction: {"V/regex/", "[global] interactively take commands to run against all lines not matching $regex"},
	globalNegSearchAction:    {"v/regex/cmd", "[global] run $cmd against each line not matching $regex"},
	superWriteAction:         {"", ""},
	writeAction:              {"w filename", "[write] write the buffer to the set filename or $filename if provided"},
	setPagerAction:           {"z n", "[pager] set pager to $n lines"},
	debugAction:              {"*", "[debug] print debugging information"},
}

var additionalDocs = map[rune]quickdoc{
	'>': {"(,)>n", "[peek] print up to n lines after the given lines"},
	'<': {"(,)<n", "[peek] print up to n lines before the given lines"},
}

func allCmds() string {
	return string(cmdOrder)
}

func excludeFromGlob(r rune) bool {
	return contains(string([]rune{
		setPagerAction,
		globalSearchAction,
		globalNegSearchAction,
		globalIntSearchAction,
		globalNegIntSearchAction,
	}), r)
}

func invertDirection(r rune) bool {
	return contains(string([]rune{
		globalNegSearchAction,
		globalNegIntSearchAction,
	}), r)
}

var cmdOrder = []rune{
	appendAction,
	breakAction,
	changeAction,
	deleteAction,
	editAction,
	reallyEditAction,
	filenameAction,
	globalSearchAction,
	globalIntSearchAction,
	helpAction,
	historyAction,
	insertAction,
	joinAction,
	copyAction,
	printLiteralAction,
	moveAction,
	mirrorAction,
	printNumsAction,
	printAction,
	quitAction,
	reallyQuitAction,
	readAction,
	shellAction,
	simpleReplaceAction,
	regexReplaceAction,
	transliterateAction,
	undoAction,
	globalNegSearchAction,
	globalNegIntSearchAction,
	writeAction,
	superWriteAction,
	setPagerAction,
	searchAction,
	searchRevAction,
	putMarkAction,
	bulkMarkAction,
	eqAction,
	debugAction,
}

func quickHelp() string {
	var str strings.Builder

	for _, r := range cmdOrder {
		if q, ok := cmds[r]; ok && len(q.example) > 0 {
			fmt.Fprintf(&str, "  - %-8s %s\n\n", q.example, wrap(q.explanation, len(q.example), 13))
		}
	}

	for r := range additionalDocs {
		if q, ok := additionalDocs[r]; ok && len(q.example) > 0 {
			fmt.Fprintf(&str, "  - %-8s %s\n\n", q.example, wrap(q.explanation, len(q.example), 13))
		}
	}
	return str.String()
}

func wrap(s string, prev, prepad int) string {
	pad := strings.Repeat(" ", prepad)
	parts := strings.Split(s, " ")

	var final strings.Builder
	var l int
	if len(s) > 80-8 || prev > 8 {
		io.WriteString(&final, "\n")
		io.WriteString(&final, pad)
	}

	for _, part := range parts {
		l += len(part) + 1
		if l > 80-13 {
			io.WriteString(&final, "\n")
			io.WriteString(&final, pad)
			l = 0
		}
		io.WriteString(&final, part)
		io.WriteString(&final, " ")
	}
	return final.String()
}
