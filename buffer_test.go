package main

import (
	"testing"
)

func getTestBuffer() buffer {
	return &memoryBuf{
		curline:  5,
		lastline: 5,
		lines: []bufferLine{
			{txt: ``, mark: false},
			{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`, mark: false},
			{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`, mark: false},
			{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`, mark: false},
			{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
			{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: false},
		},
		filename: "filename",
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
		{"1", "7", 5, 5, 1, 7, false},
		{"0", "7", 5, 5, 0, 0, true},
	}

	for _, test := range tests {
		controlBuffer := getTestBuffer()
		l1, l2, err := controlBuffer.defLines(test.start, test.end, test.l1, test.l2)
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
