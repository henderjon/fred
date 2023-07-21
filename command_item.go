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
		return fmt.Sprintf("type: %s; val: %s", itemTyps[i.typ], i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("type: %s; val: %.10q", itemTyps[i.typ], i.val)
	}
	return fmt.Sprintf("type: %s; val: %q", itemTyps[i.typ], i.val)
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
	itemSubstitution
	itemAdditional
	itemEmpty
	itemCommand
	itemUnknownCommand
	itemEmptyPattern
	itemMissingDelim
	itemGlobalFlag
)

var itemTyps = map[itemType]string{
	itemError:          "itemError",
	itemEOF:            "itemEOF",
	itemKeyword:        "itemKeyword",
	itemNumber:         "itemNumber", // a single address
	itemRange:          "itemRange",  // a range of addresses
	itemAction:         "itemAction",
	itemDelim:          "itemDelim",
	itemPattern:        "itemPattern",
	itemSubstitution:   "itemSubstitution",
	itemAdditional:     "itemAdditional",
	itemEmpty:          "itemEmpty",
	itemCommand:        "itemCommand",
	itemUnknownCommand: "itemUnknownCommand",
	itemEmptyPattern:   "itemEmptyPattern",
	itemMissingDelim:   "itemMissingDelim",
	itemGlobalFlag:     "itemGlobalFlag",
}
