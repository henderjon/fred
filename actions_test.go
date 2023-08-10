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
			{txt: ``},
			{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
			{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
			{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
			{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
			{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
		},
		filename: "filename",
		dirty:    false,
	}
}

func getTestMarkedActionBuffer() buffer {
	return &memoryBuf{
		curline:  1,
		lastline: 5,
		lines: []bufferLine{
			{txt: ``},
			{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`, mark: mark},
			{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`, mark: mark},
			{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`, mark: mark},
			{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: mark},
			{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`, mark: mark},
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
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
			},
			filename: "filename",
			dirty:    true,
		}},
		{1, 5, &memoryBuf{
			curline:  0,
			lastline: 0,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
		}},
	}

	for i, test := range tests {
		controlBuffer := getTestActionBuffer()
		doDelete(controlBuffer, test.l1, test.l2)

		if diff := cmp.Diff(controlBuffer, test.expected, cmp.AllowUnexported(memoryBuf{}, bufferLine{}, search{})); diff != "" {
			t.Errorf("idx: %d; -got/+want\n%s", i, diff)
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
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
			},
			filename: "filename",
			dirty:    true,
		}},
		{4, 5, "0", &memoryBuf{
			curline:  2,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
			},
			filename: "filename",
			dirty:    true,
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
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
			},
			filename: "filename",
			dirty:    true,
		}},
		{4, 5, "0", &memoryBuf{
			curline:  0,
			lastline: 7,
			lines: []bufferLine{
				{txt: ``},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
		}},
	}

	for i, test := range tests {
		controlBuffer := getTestActionBuffer()
		err := doCopyNPaste(controlBuffer, test.l1, test.l2, test.dest)
		if err != nil {
			t.Error(err)
		}

		if diff := cmp.Diff(controlBuffer, test.expected, cmp.AllowUnexported(memoryBuf{}, bufferLine{}, search{})); diff != "" {
			t.Errorf("idx: %d; -got/+want\n%s", i, diff)
		}

	}
}
func Test_doSimpleReplace(t *testing.T) {
	tests := []struct {
		l1, l2   int
		pattern  string
		replace  string
		num      string
		expected buffer
	}{
		{2, 4, "or", "-", "1", &memoryBuf{
			curline:  4,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut p-ta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida p-ttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
		}},
		{1, 5, "or", "-", "3", &memoryBuf{
			curline:  5,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. M-bi sed ante eu ...`},
				{txt: `2 Duis ut porta mi, eu ornare -ci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
		}},
		{2, 4, "or", "-", "-1", &memoryBuf{
			curline:  4,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut p-ta mi, eu -nare -ci. Etiam sed vehicula -ci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida p-ttit-. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
		}},
	}

	for _, test := range tests {
		controlBuffer := getTestActionBuffer()
		doSimpleReplace(controlBuffer, test.l1, test.l2, test.pattern, test.replace, test.num, &cache{})

		if diff := cmp.Diff(controlBuffer, test.expected, cmp.AllowUnexported(memoryBuf{}, bufferLine{}, search{})); diff != "" {
			t.Errorf("-got/+want\n%s", diff)
		}

	}
}

func Test_doRegexReplace(t *testing.T) {
	tests := []struct {
		l1, l2   int
		pattern  string
		replace  string
		num      string
		expected buffer
	}{
		{2, 4, "or", "-", "1", &memoryBuf{
			curline:  4,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut p-ta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida p-ttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
		}},
		{1, 5, "or", "-", "3", &memoryBuf{
			curline:  5,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. M-bi sed ante eu ...`},
				{txt: `2 Duis ut porta mi, eu ornare -ci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
		}},
		{2, 4, "or", "-", "-1", &memoryBuf{
			curline:  4,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut p-ta mi, eu -nare -ci. Etiam sed vehicula -ci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida p-ttit-. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
		}},
	}

	for _, test := range tests {
		controlBuffer := getTestActionBuffer()
		doRegexReplace(controlBuffer, test.l1, test.l2, test.pattern, test.replace, test.num, &cache{})

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
			globalPrefix: "g",
		}, &memoryBuf{
			curline:  2,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 L%em ipsum dol% sit amet, consectetur adipiscing elit. M%bi sed ante eu ...`},
				{txt: `2 Duis ut p%ta mi, eu %nare %ci. Etiam sed vehicula %ci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida p%ttit%. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
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
			globalPrefix: "g",
		}, &memoryBuf{
			curline:  2,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 L%em ipsum dol% sit amet, consectetur adipiscing elit. M%bi sed ante eu ...`},
				{txt: `2 Duis ut p%ta mi, eu %nare %ci. Etiam sed vehicula %ci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
		}},
	}

	for i, test := range tests {
		input := getInput(os.Stdin, os.Stdout)
		controlBuffer := getTestActionBuffer()
		doGlobMarks(test.cmd, controlBuffer)
		doGlob(test.cmd, controlBuffer, input, &cache{})

		if diff := cmp.Diff(controlBuffer, test.expected, cmp.AllowUnexported(memoryBuf{}, bufferLine{}, search{})); diff != "" {
			t.Errorf("idx: %d; -got/+want\n%s", i, diff)
		}

	}
}

func Test_doSetMarkLine_1(t *testing.T) {
	tests := []struct {
		l1, l2   int
		argument string
		expected buffer
	}{
		{2, 4, "foo", &memoryBuf{
			curline:  1,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`, mark: 'f'},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`, mark: 'f'},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`, mark: 'f'},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
		}},
	}

	for _, test := range tests {
		controlBuffer := getTestActionBuffer()
		doSetMarkLine(controlBuffer, test.l1, test.l2, test.argument)

		if diff := cmp.Diff(controlBuffer, test.expected, cmp.AllowUnexported(memoryBuf{}, bufferLine{}, search{})); diff != "" {
			t.Errorf("-got/+want\n%s", diff)
		}

	}
}

func Test_doSetMarkLine_2(t *testing.T) {
	tests := []struct {
		l1, l2   int
		argument string
		expected buffer
	}{
		{4, 5, "", &memoryBuf{
			curline:  1,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`, mark: mark},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`, mark: mark},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`, mark: mark},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
		}},
	}

	for _, test := range tests {
		controlBuffer := getTestMarkedActionBuffer()
		doSetMarkLine(controlBuffer, test.l1, test.l2, test.argument)

		if diff := cmp.Diff(controlBuffer, test.expected, cmp.AllowUnexported(memoryBuf{}, bufferLine{}, search{})); diff != "" {
			t.Errorf("-got/+want\n%s", diff)
		}

	}
}
