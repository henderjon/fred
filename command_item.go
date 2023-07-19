package main

import "fmt"

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType
	val string
}

func (i *item) String() string {
	switch {
	case i.typ == itemEOF:
		// return "EOF"
		fallthrough
	case i.typ == itemError:
		// return i.val
		fallthrough
	case i.typ > itemKeyword:
		return fmt.Sprintf("%s <%s>", itemTyps[i.typ], i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("%s %.10q...", itemTyps[i.typ], i.val)
	}
	return fmt.Sprintf("%s %q", itemTyps[i.typ], i.val)
}

// itemType identifies the type of lex items.
type itemType int

// [num][range][num][action][delim][pattern][delim][pattern][delim][additional]

const (
	itemError itemType = iota + 1 // error occurred; value is text of error
	itemEOF
	itemKeyword // used only to delimit the keywords
	itemNumber
	itemRange
	itemAction
	itemDelim
	itemPattern
	itemAdditional
	itemEmpty
)

var itemTyps = map[itemType]string{
	itemError:      "itemError",
	itemEOF:        "itemEOF",
	itemKeyword:    "itemKeyword",
	itemNumber:     "itemNumber",
	itemRange:      "itemRange",
	itemAction:     "itemAction",
	itemDelim:      "itemDelim",
	itemPattern:    "itemPattern",
	itemAdditional: "itemAdditional",
	itemEmpty:      "itemEmpty",
}