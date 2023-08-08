package main

import "errors"

const (
	zero = "1"
	end  = "$"
)

type parser struct{}

func (p *parser) run(input string) (*command, error) {
	var (
		i item
		c = &command{}
	)

	// stack := make(chan item, 2)
	l := lex(input, "main input")
	go l.run()

	// Address: break Address
	for {
		i = l.nextItem()
		switch i.typ {
		case itemRange: // did we get a comma first?
			c.setAddr(zero)
			c.setAddr(end)
		case itemAddress: // did we get a number first?
			c.setAddr(i.val)
		case itemAction: // did we get an action first
			c.setAction(rune(i.val[0]))
		case itemDestination:
			c.setDestination(i.val)
		// case itemAdditional:
		// c.setSubCommand(i.val)
		case itemAddressPattern:
			c.setAddrPattern(i.val)
		case itemPattern:
			c.setPattern(i.val)
		case itemGlobalPrefix:
			c.setGlobalPrefix(i.val)
		case itemReplaceNum: // this takes the
			c.setReplaceNum(i.val)
		case itemSubstitution:
			c.setSubstitution(i.val)
		case itemArg:
			c.setArgument(i.val)
		case itemError:
			return nil, errors.New(i.val)
		case itemEOF: // no more items
			return c, nil
		default:
			stderr.Log(i.String())
			return c, nil
		}
	}
	// return c, nil
}
