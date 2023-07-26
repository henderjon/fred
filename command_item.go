package main

import "fmt"

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType
	val string
}

func (i *item) String() string {
	return fmt.Sprintf("type: %s; val: %q", itemTypes[i.typ], i.val)
}

// itemType identifies the type of lex items.
type itemType int

// [num][range][num][action][delim][pattern][delim][pattern][delim][additional]

const (
	itemError itemType = iota + 1 // error occurred; value is text of error
	itemEOF
	itemAddress
	itemAddressPattern
	itemRange
	itemAction
	itemPattern
	itemSubstitution
	itemDestination
	itemAdditional
	itemGlobalPrefix
	itemReplaceNum
)

var itemTypes = map[itemType]string{
	itemError:          "itemError",
	itemEOF:            "itemEOF",
	itemAddress:        "itemAddress",        // a single address
	itemAddressPattern: "itemAddressPattern", // find lines matching this pattern
	itemRange:          "itemRange",          // a range of addresses
	itemAction:         "itemAction",
	itemPattern:        "itemPattern",
	itemSubstitution:   "itemSubstitution",
	itemDestination:    "itemDestination",
	itemAdditional:     "itemAdditional",
	itemGlobalPrefix:   "itemGlobalPrefix",
	itemReplaceNum:     "itemReplaceNum",
}
