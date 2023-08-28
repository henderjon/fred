//go:build memory

package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func getTestBuffer() buffer {
	return &memoryBuf{
		curline:  5,
		lastline: 5,
		lines: []bufferLine{
			{txt: ``},
			{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
			{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
			{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
			{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
			{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
		},
		filename: "filename",
	}
}

func Test_sameBuffer(t *testing.T) {
	one := getTestBuffer()
	two := getTestBuffer()

	// is it possible to compare two buffers?
	if one == two {
		t.Error("buffers the same", one, two)
	}
}

func Test_defLines(t *testing.T) {
	tests := []struct {
		start, end string
		l1, l2     int
		i1, i2     int
		e          bool
	}{
		{"2", "4", 5, 5, 2, 4, false},
		{".", "$", 5, 5, 5, 5, false},
		{"1", "$", 5, 5, 1, 5, false}, // ',' injects these values
		{"1", "7", 5, 5, 1, 5, false},
		{"0", "7", 5, 5, 1, 5, false}, // we coerce out of bounds addresses, now
	}

	for _, test := range tests {
		controlBuffer := getTestBuffer()
		l1, l2, err := controlBuffer.defLines(test.start, test.end, "", test.l1, test.l2)
		if err != nil && test.e == false {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}

		if err == nil && test.e == true {
			t.Errorf("expected error; given (%s,%s)", test.start, test.end)
			continue
		}

		if l1 != test.i1 {
			t.Errorf("error converting first address; given (%s,%s); got %d; want %d", test.start, test.end, l1, test.i1)
		}

		if l2 != test.i2 {
			t.Errorf("error converting second address; given (%s,%s); got %d; want %d", test.start, test.end, l2, test.i2)
		}
	}
}
func Test_guardAddresss(t *testing.T) {
	tests := []struct {
		addr     string
		expected int
	}{
		{"2", 2},
		{".", 5},
		{"", 5},
	}

	for _, test := range tests {
		controlBuffer := getTestBuffer()
		i, err := guardAddress(controlBuffer, test.addr)
		if err != nil {
			t.Error(err)
		}

		if i != test.expected {
			t.Errorf("error converting first address; given '%s'; got %d; want %d", test.addr, i, test.expected)
		}
	}
}

func Test_scan(t *testing.T) {
	controlBuffer := getTestBuffer()
	tests := []struct {
		label    string
		scanner  func() (int, bool)
		expected []int
	}{
		{"forward: 1, 5-1", controlBuffer.scanForward(1, 5-1), []int{1, 2, 3, 4, 5}},
		{"forward: 3, 5-3", controlBuffer.scanForward(3, 5-3), []int{3, 4, 5}},
		{"forward: 2, 4-2", controlBuffer.scanForward(2, 4-2), []int{2, 3, 4}},
		{"forward: 4, 5", controlBuffer.scanForward(4, 5), []int{4, 5, 0, 1, 2, 3}},
		{"forward: 1, 0", controlBuffer.scanForward(1, 0), []int{1}},
		{"reverse: 5, 5-1", controlBuffer.scanReverse(1, 5-1), []int{1, 0, 5, 4, 3}},
		{"reverse: 3, 5-3", controlBuffer.scanReverse(3, 5-3), []int{3, 2, 1}},
		{"reverse: 2, 4-2", controlBuffer.scanReverse(2, 4-2), []int{2, 1, 0}},
		{"reverse: 2, 2-4", controlBuffer.scanReverse(2, 4-2), []int{2, 1, 0}},
		{"reverse: 4, 5", controlBuffer.scanReverse(4, 5), []int{4, 3, 2, 1, 0, 5}},
	}

	for _, test := range tests {
		results := make([]int, 0)
		for {
			i, ok := test.scanner()
			if !ok {
				break
			}
			results = append(results, i)
		}

		if diff := cmp.Diff(results, test.expected); diff != "" {
			t.Errorf("given: %s; -got/+want\n%s", test.label, diff)
		}
	}
}

func Test_simpleNReplace(t *testing.T) {
	tests := []struct {
		subject  string
		pattern  string
		replace  string
		n        int
		expected string
	}{
		{"one one one one one", "one", "two", 3, "one one two one one"},
		{"one one one one one", "one", "two", 1, "two one one one one"},
		{"one one one one one", "one", "two", 6, "one one one one one"},
		{"one one one one one", "six", "two", 6, "one one one one one"},
	}

	for _, test := range tests {
		result := simpleNReplace(test.subject, test.pattern, test.replace, test.n)

		if diff := cmp.Diff(result, test.expected); diff != "" {
			t.Errorf("given: %s; -got/+want\n%s", test.subject, diff)
		}

	}
}

func Test_handleTabs(t *testing.T) {
	tests := []struct {
		before string
		after  string
	}{
		{`asdf\\t\tasdf\\t\tasdf`, `asdf\t	asdf\t	asdf`},
		{`\\t	\t	\\t`, `\t			\t`},
		{`\\t\t\\t`, `\t	\t`},
	}

	for _, test := range tests {
		result := handleTabs(test.before)

		if diff := cmp.Diff(result, test.after); diff != "" {
			t.Errorf("given: %s; -got/+want\n%s", test.before, diff)
		}

	}
}
