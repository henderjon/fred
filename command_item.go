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
		return fmt.Sprintf("type: %s; val: %s", itemTypes[i.typ], i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("type: %s; val: %.10q", itemTypes[i.typ], i.val)
	}
	return fmt.Sprintf("type: %s; val: %q", itemTypes[i.typ], i.val)
}

// itemType identifies the type of lex items.
type itemType int

// [num][range][num][action][delim][pattern][delim][pattern][delim][additional]

const (
	itemError itemType = iota + 1 // error occurred; value is text of error
	itemEOF
	itemKeyword
	itemAddress
	itemRange
	itemAction
	itemDelim
	itemPattern
	itemSubstitution
	itemDestination
	itemCommand
	itemGlobalPrefix
)

var itemTypes = map[itemType]string{
	itemError:        "itemError",
	itemEOF:          "itemEOF",
	itemKeyword:      "itemKeyword",
	itemAddress:      "itemAddress", // a single address
	itemRange:        "itemRange",   // a range of addresses
	itemAction:       "itemAction",
	itemDelim:        "itemDelim",
	itemPattern:      "itemPattern",
	itemSubstitution: "itemSubstitution",
	itemDestination:  "itemDestination",
	itemCommand:      "itemCommand",
	itemGlobalPrefix: "itemGlobalPrefix",
}
