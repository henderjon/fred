package main

type parser struct{}

func (p *parser) run(input string) *command {
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
		case itemAction: // did we get an action first
			c.action = i.val
			break Address
		case itemEOF: // no more items
			return c // skip the rest of the func
		default:
			stderr.Log(i.String())
			stderr.Log("first token should be of type(s): action, range, number")
		}
	}
	return c
}
