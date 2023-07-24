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
// 		{"12", &command{addRange: []int{12}}},
// 		{"-12", &command{addRange: []int{0}}},
// 		{"+12", &command{addRange: []int{12}}},
// 		{"   12  ", &command{addRange: []int{12}}},
// 		{"   12,  ", &command{addRange: []int{12, -1}}},         // parses correctly but invalid address?
// 		{"   12,13  ", &command{addRange: []int{12, 13}}},       // parses correctly but invalid address?
// 		{"   12,13,14,15  ", &command{addRange: []int{12, 15}}}, // parses correctly but invalid address?
// 		{"   12,13,14,11  ", &command{addRange: []int{12, 12}}}, // parses correctly but invalid address?
// 		{"  ,12  ", &command{addRange: []int{0, 12}}},
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
		// typical usage
		{",n", &command{addRange: []int{0, -1}, action: printNumsAction}, false},
		{"10,25p", &command{addRange: []int{10, 25}, action: printAction}, false},
		{"10,25n", &command{addRange: []int{10, 25}, action: printNumsAction}, false},
		{"q", &command{action: quitAction}, false},
		{"25a", &command{addRange: []int{25}, action: appendAction}, false},
		{"10,25d", &command{addRange: []int{10, 25}, action: deleteAction}, false},
		{"10,25i", &command{addRange: []int{10, 25}, action: insertAction}, false},
		{"10,25c", &command{addRange: []int{10, 25}, action: changeAction}, false},
		{"=", &command{action: eqAction}, false},
		{"10,25m35", &command{addRange: []int{10, 25}, action: moveAction, destination: 35}, false},
		{"10,25k35", &command{addRange: []int{10, 25}, action: copyAction, destination: 35}, false},
		// testing edge cases
		{"12", &command{addRange: []int{12}}, false},
		{"-12", &command{addRange: []int{0}}, false},
		{"+12", &command{addRange: []int{12}}, false},
		{"   12  ", &command{addRange: []int{12}}, false},
		{"   12,  ", &command{addRange: []int{12, -1}}, false},
		{"   12,13  ", &command{addRange: []int{12, 13}}, false},
		{"   12,13,14,15  ", &command{addRange: []int{12, 15}}, false}, // parses correctly but invalid address?
		{"   12,13,14,11  ", &command{addRange: []int{12, 12}}, false}, // parses correctly but invalid address?
		{"  ,12  ", &command{addRange: []int{0, 12}}, false},
		{"12a", &command{addRange: []int{12}, action: appendAction}, false},
		{"12,230a", &command{addRange: []int{12, 230}, action: appendAction}, false},
		{"+12i", &command{addRange: []int{12}, action: insertAction}, false},
		{"   12 a  ", &command{addRange: []int{12}, action: appendAction}, false},
		{"   12,a  ", &command{addRange: []int{12, -1}, action: appendAction}, false},
		{"   12 b  ", nil, true},                                                  // unknown command
		{"g/^f[ob]ar/", &command{globalPrefix: true, pattern: `^f[ob]ar`}, false}, // missing address
		{",g/^f[ob]ar/", &command{addRange: []int{0, -1}, globalPrefix: true, pattern: `^f[ob]ar`}, false},
		{"5,g/^f[ob]ar/", &command{addRange: []int{5, -1}, globalPrefix: true, pattern: `^f[ob]ar`}, false},
		{"5,8g/^f[ob]ar/", &command{addRange: []int{5, 8}, globalPrefix: true, pattern: `^f[ob]ar`}, false},
		{"5,8g/^f[ob]ar/p", &command{addRange: []int{5, 8}, action: printAction, globalPrefix: true, pattern: `^f[ob]ar`}, false},
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

// {"10,25/"} searchAction     = rune('/') // /re/... establishes the ADDRESSES for the lines against which to execute cmd

// {"10,25g"} gSearchAction    = rune('g') // g/re/p is the glob prefix which means we use the pattern to print every line that matches [gPREFIX]
// {"10,25s"} substituteAction = rune('s') // Substitute CMD
// {"10,25e"} editAction       = rune('e') // Edit command
// {"10,25f"} fileAction       = rune('f') // File command
// {"10,25r"} readAction       = rune('r') // read [file] command
// {"10,25w"} writeAction      = rune('w') // write [file] command
