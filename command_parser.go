package main

import "errors"

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
			// c.setAddr(i.val) // TODO: this should be unnecessary
		case itemAction: // did we get an action first
			c.setAction(rune(i.val[0]))
		case itemDestination:
			c.setDestination(i.val)
		case itemPattern: // no more items
			c.setPattern(i.val)
		case itemGlobalPrefix: // no more items
			c.setGlobalPrefix(true)
		case itemSubstitution: // no more items
			c.setSubstitution(i.val)
		case itemError: // no more items
			return nil, errors.New(i.val) // skip the rest of the func
		case itemEOF: // no more items
			return c, nil // skip the rest of the func
		default:
			stderr.Log(i.String())
			return c, nil
		}
	}
	// return c, nil
}
