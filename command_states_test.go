package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_lexNumber(t *testing.T) {
	// we're testing our accuracy in parsing NUMBERS, parsing the full address(es) for the command is tested in the parser
	tests := []struct {
		input    string
		expected item
	}{
		{"12", item{itemNumber, "12"}},
		{"-12", item{itemNumber, "-12"}},
		{"+12", item{itemNumber, "+12"}},
		{"   12  ", item{itemNumber, "12"}},
		{"   12,  ", item{itemNumber, "12"}},
		{"  ,12  ", item{itemRange, ","}},
	} //itemEmpty
	var i item
	for _, test := range tests {
		l := lex("", test.input)
		go l.run()
		i = l.nextItem()
		if diff := cmp.Diff(i, test.expected, cmp.AllowUnexported(item{})); diff != "" {
			t.Errorf("%s; (from) -got/+want\n%s", test.input, diff)
		}
	}
}

func Test_lexCommand(t *testing.T) {
	// we're testing our accuracy in parsing NUMBERS, parsing the full address(es) for the command is tested in the parser
	tests := []struct {
		input    string
		expected item
	}{
		{"a", item{itemAction, string(appendAction)}},
		{"i", item{itemAction, string(insertAction)}},
		{"c", item{itemAction, string(changeAction)}},
		{"b", item{itemUnknownCommand, "b"}},
	} //itemEmpty
	var i item
	for _, test := range tests {
		l := lex("", test.input)
		go l.run()
		i = l.nextItem()
		// stderr.Fatal(i.String())
		if diff := cmp.Diff(i, test.expected, cmp.AllowUnexported(item{})); diff != "" {
			t.Errorf("%s; (from) -got/+want\n%s", test.input, diff)
		}
	}
}
