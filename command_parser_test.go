package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// func Test_parser_run(t *testing.T) {
// 	tests := []struct {
// 		input    string
// 		expected *command
// 	}{
// 		{"12", &command{addrs: []int{12}}},
// 		{"-12", &command{addrs: []int{0}}},
// 		{"+12", &command{addrs: []int{12}}},
// 		{"   12  ", &command{addrs: []int{12}}},
// 		{"   12,  ", &command{addrs: []int{12, -1}}},         // parses correctly but invalid address?
// 		{"   12,13  ", &command{addrs: []int{12, 13}}},       // parses correctly but invalid address?
// 		{"   12,13,14,15  ", &command{addrs: []int{12, 15}}}, // parses correctly but invalid address?
// 		{"   12,13,14,11  ", &command{addrs: []int{12, 12}}}, // parses correctly but invalid address?
// 		{"  ,12  ", &command{addrs: []int{0, 12}}},
// 	} //itemEmpty

// 	for _, test := range tests {
// 		c, _ := (&parser{}).run(test.input)

// 		if diff := cmp.Diff(c, test.expected, cmp.AllowUnexported(command{})); diff != "" {
// 			t.Errorf("given: %s; -got/+want\n%s", test.input, diff)
// 		}
// 	}
// }

func Test_parser_full_commands(t *testing.T) {
	tests := []struct {
		input      string
		expCommand *command
		expErr     bool
	}{
		{"12", &command{addrs: []int{12}}, false},
		{"-12", &command{addrs: []int{0}}, false},
		{"+12", &command{addrs: []int{12}}, false},
		{"   12  ", &command{addrs: []int{12}}, false},
		{"   12,  ", &command{addrs: []int{12, -1}}, false},
		{"   12,13  ", &command{addrs: []int{12, 13}}, false},
		{"   12,13,14,15  ", &command{addrs: []int{12, 15}}, false}, // parses correctly but invalid address?
		{"   12,13,14,11  ", &command{addrs: []int{12, 12}}, false}, // parses correctly but invalid address?
		{"  ,12  ", &command{addrs: []int{0, 12}}, false},
		{"12a", &command{addrs: []int{12}, action: appendAction}, false},
		{"12,230a", &command{addrs: []int{12, 230}, action: appendAction}, false},
		{"+12i", &command{addrs: []int{12}, action: insertAction}, false},
		{"   12 a  ", &command{addrs: []int{12}, action: appendAction}, false},
		{"   12,a  ", &command{addrs: []int{12, -1}, action: appendAction}, false},
		{"   12 b  ", nil, true},                                                  // unknown command
		{"g/^f[ob]ar/", &command{globalPrefix: true, pattern: `^f[ob]ar`}, false}, // missing address
		{",g/^f[ob]ar/", &command{addrs: []int{0, -1}, globalPrefix: true, pattern: `^f[ob]ar`}, false},
		{"5,g/^f[ob]ar/", &command{addrs: []int{5, -1}, globalPrefix: true, pattern: `^f[ob]ar`}, false},
		{"5,8g/^f[ob]ar/", &command{addrs: []int{5, 8}, globalPrefix: true, pattern: `^f[ob]ar`}, false},
		{"5,8g/^f[ob]ar/p", &command{addrs: []int{5, 8}, action: printAction, globalPrefix: true, pattern: `^f[ob]ar`}, false},
		{"g/^f[ob]ar/", &command{globalPrefix: true, pattern: `^f[ob]ar`}, false},
		{"g/b.g/", &command{globalPrefix: true, pattern: `b.g`}, false},
		{"g//", nil, true},   //itemEmptyPattern
		{"g/b.z", nil, true}, //itemMissingDelim
		{"/re/p", &command{action: printAction, pattern: `re`}, false},
	} //itemEmpty

	for _, test := range tests {
		c, err := (&parser{}).run(test.input)

		if diff := cmp.Diff(c, test.expCommand, cmp.AllowUnexported(command{})); diff != "" {
			t.Errorf("given: %s; -got/+want\n%s[%s]", test.input, diff, err)
			// t.Errorf("given: %s; -got/+want\n%s", test.input, diff)
		}

		if test.expErr && (err == nil) {
			t.Errorf("given: %s; error: %s\n", test.input, err)
		}

	}
}
