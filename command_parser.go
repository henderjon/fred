package main

import "errors"

type parser struct{}

func (p *parser) run(input string) (*command, error) {
	var (
		i item
		c = &command{}
	)

	l := lex("", input)
	go l.run()

Address:
	for {
		i = l.nextItem()
		switch i.typ {
		case itemRange: // did we get a comma first?
			c.setFrom(zero)
			c.setFrom(end)
		case itemNumber: // did we get a number first?
			c.setFrom(i.val)
			// c.setFrom(i.val) // TODO: this should be unnecessary
		case itemAction: // did we get an action first
			c.setAction(rune(i.val[0]))
			break Address
		case itemUnknownCommand: // no more items
			return c, errors.New("unknonwn command") // skip the rest of the func
		case itemError: // no more items
			return c, errors.New("unknonwn error") // skip the rest of the func
		case itemEOF: // no more items
			return c, errors.New("eof") // skip the rest of the func
		default:
			stderr.Log("first token should be of type(s): action, range, number")
			return c, nil
		}
	}
	return c, nil
}
