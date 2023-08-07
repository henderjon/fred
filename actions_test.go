package main

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func getTestActionBuffer() buffer {
	return &memoryBuf{
		curline:  1,
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

func Test_doDelete(t *testing.T) {
	tests := []struct {
		l1, l2   int
		expected buffer
	}{
		{2, 4, &memoryBuf{
			curline:  1,
			lastline: 2,
			lines: []bufferLine{
				{txt: ``, mark: false},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`, mark: false},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: false},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`, mark: false},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`, mark: false},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
			},
			filename: "filename",
		}},
		{1, 5, &memoryBuf{
			curline:  0,
			lastline: 0,
			lines: []bufferLine{
				{txt: ``, mark: false},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`, mark: false},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`, mark: false},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`, mark: false},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: false},
			},
			filename: "filename",
		}},
	}

	for _, test := range tests {
		controlBuffer := getTestActionBuffer()
		doDelete(controlBuffer, test.l1, test.l2)

		if diff := cmp.Diff(controlBuffer, test.expected, cmp.AllowUnexported(memoryBuf{}, bufferLine{}, search{})); diff != "" {
			t.Errorf("-got/+want\n%s", diff)
		}

	}
}
func Test_doMove(t *testing.T) {
	tests := []struct {
		l1, l2   int
		dest     string
		expected buffer
	}{
		{2, 4, "5", &memoryBuf{
			curline:  5,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``, mark: false},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`, mark: false},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: false},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`, mark: false},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`, mark: false},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
			},
			filename: "filename",
		}},
		{4, 5, "0", &memoryBuf{
			curline:  2,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``, mark: false},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: false},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`, mark: false},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`, mark: false},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`, mark: false},
			},
			filename: "filename",
		}},
	}

	for _, test := range tests {
		controlBuffer := getTestActionBuffer()
		doMove(controlBuffer, test.l1, test.l2, test.dest)

		if diff := cmp.Diff(controlBuffer, test.expected, cmp.AllowUnexported(memoryBuf{}, bufferLine{}, search{})); diff != "" {
			t.Errorf("-got/+want\n%s", diff)
		}

	}
}

func Test_doCopyNPaste(t *testing.T) {
	tests := []struct {
		l1, l2   int
		dest     string
		expected buffer
	}{
		{2, 4, "5", &memoryBuf{
			curline:  5,
			lastline: 8,
			lines: []bufferLine{
				{txt: ``, mark: false},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`, mark: false},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`, mark: false},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`, mark: false},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: false},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`, mark: false},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`, mark: false},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
			},
			filename: "filename",
		}},
		{4, 5, "0", &memoryBuf{
			curline:  0,
			lastline: 7,
			lines: []bufferLine{
				{txt: ``, mark: false},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: false},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`, mark: false},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`, mark: false},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`, mark: false},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: false},
			},
			filename: "filename",
		}},
	}

	for _, test := range tests {
		controlBuffer := getTestActionBuffer()
		doCopyNPaste(controlBuffer, test.l1, test.l2, test.dest)

		if diff := cmp.Diff(controlBuffer, test.expected, cmp.AllowUnexported(memoryBuf{}, bufferLine{}, search{})); diff != "" {
			t.Errorf("-got/+want\n%s", diff)
		}

	}
}
func Test_doSimpleReplace(t *testing.T) {
	tests := []struct {
		l1, l2   int
		num      string
		expected buffer
	}{
		{2, 4, "1", &memoryBuf{
			curline:  4,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``, mark: false},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`, mark: false},
				{txt: `2 Duis ut p-ta mi, eu ornare orci. Etiam sed vehicula orci. ...`, mark: false},
				{txt: `3 Nunc scelerisque urna a erat gravida p-ttitor. Donec pulvinar leo urna, id ...`, mark: false},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: false},
			},
			filename: "filename",
		}},
		{2, 4, "-1", &memoryBuf{
			curline:  4,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``, mark: false},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`, mark: false},
				{txt: `2 Duis ut p-ta mi, eu -nare -ci. Etiam sed vehicula -ci. ...`, mark: false},
				{txt: `3 Nunc scelerisque urna a erat gravida p-ttit-. Donec pulvinar leo urna, id ...`, mark: false},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: false},
			},
			filename: "filename",
		}},
	}

	for _, test := range tests {
		controlBuffer := getTestActionBuffer()
		doSimpleReplace(controlBuffer, test.l1, test.l2, "or", "-", test.num)

		if diff := cmp.Diff(controlBuffer, test.expected, cmp.AllowUnexported(memoryBuf{}, bufferLine{}, search{})); diff != "" {
			t.Errorf("-got/+want\n%s", diff)
		}

	}
}
func Test_doGlob(t *testing.T) {
	tests := []struct {
		cmd      command
		expected buffer
	}{
		{command{
			// addrStart:    "",
			// addrEnd:      "",
			addrPattern:  "[1-5]",
			action:       simpleReplaceAction,
			pattern:      "or",
			substitution: "%",
			replaceNum:   "-1",
			// destination:  "",
			// subCommand:   "",
			// argument:     "",
			globalPrefix: true,
		}, &memoryBuf{
			curline:  5,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``, mark: false},
				{txt: `1 L%em ipsum dol% sit amet, consectetur adipiscing elit. M%bi sed ante eu ...`, mark: false},
				{txt: `2 Duis ut p%ta mi, eu %nare %ci. Etiam sed vehicula %ci. ...`, mark: false},
				{txt: `3 Nunc scelerisque urna a erat gravida p%ttit%. Donec pulvinar leo urna, id ...`, mark: false},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: false},
			},
			filename: "filename",
		}},
		{command{
			addrStart:    "1",
			addrEnd:      "2",
			addrPattern:  "[1-5]",
			action:       simpleReplaceAction,
			pattern:      "or",
			substitution: "%",
			replaceNum:   "-1",
			// destination:  "",
			// subCommand:   "",
			// argument:     "",
			globalPrefix: true,
		}, &memoryBuf{
			curline:  2,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``, mark: false},
				{txt: `1 L%em ipsum dol% sit amet, consectetur adipiscing elit. M%bi sed ante eu ...`, mark: false},
				{txt: `2 Duis ut p%ta mi, eu %nare %ci. Etiam sed vehicula %ci. ...`, mark: false},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`, mark: false},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: false},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: false},
			},
			filename: "filename",
		}},
	}

	for _, test := range tests {
		input := getInput(os.Stdin, os.Stdout)
		controlBuffer := getTestActionBuffer()
		doGlob(test.cmd, controlBuffer, input)

		if diff := cmp.Diff(controlBuffer, test.expected, cmp.AllowUnexported(memoryBuf{}, bufferLine{}, search{})); diff != "" {
			t.Errorf("-got/+want\n%s", diff)
		}

	}
}
