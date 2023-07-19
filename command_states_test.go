package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_lexNumber(t *testing.T) {
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
