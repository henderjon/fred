package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parser_run(t *testing.T) {
	tests := []struct {
		input    string
		expected *command
	}{
		{"12", &command{addrs: []int{12}}},
		{"-12", &command{addrs: []int{-12}}},
		{"+12", &command{addrs: []int{12}}},
		{"   12  ", &command{addrs: []int{12}}},
		{"   12,  ", &command{addrs: []int{12, -1}}},         // parses correctly but invalid address?
		{"   12,13  ", &command{addrs: []int{12, 13}}},       // parses correctly but invalid address?
		{"   12,13,14,15  ", &command{addrs: []int{12, 15}}}, // parses correctly but invalid address?
		{"   12,13,14,11  ", &command{addrs: []int{12, 12}}}, // parses correctly but invalid address?
		{"  ,12  ", &command{addrs: []int{0, 12}}},
	} //itemEmpty

	for _, test := range tests {
		c := (&parser{}).run(test.input)

		if diff := cmp.Diff(c, test.expected, cmp.AllowUnexported(command{})); diff != "" {
			t.Errorf("%s; (from) -got/+want\n%s", test.input, diff)
		}
	}
}
