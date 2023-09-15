//go:build memory

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
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
			rev:      4,
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
			rev:      1,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doDelete ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			doDelete(controlBuffer, test.l1, test.l2)

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
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
			rev:      4,
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
			rev:      4,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doMove ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			doMove(controlBuffer, test.l1, test.l2, test.dest)

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
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
			rev:      4,
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
			rev:      6,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doCopyNPaste ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			err := doCopyNPaste(controlBuffer, test.l1, test.l2, test.dest)
			if err != nil {
				t.Error(err)
			}

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
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
			rev:      3,
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
			rev:      5,
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
			rev:      3,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doSimpleReplace ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			c := &cache{}
			doSimpleReplace(controlBuffer, test.l1, test.l2, c.replace(test.pattern, test.replace, test.num))

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
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
			rev:      3,
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
			rev:      5,
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
			rev:      3,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doRegexReplace ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			c := &cache{}
			doRegexReplace(controlBuffer, test.l1, test.l2, c.replace(test.pattern, test.replace, test.num))

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}
func Test_doGlob(t *testing.T) {
	tests := []struct {
		line1, line2 int
		cmd          command
		expected     buffer
	}{
		{1, 5, command{
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
			globalPrefix: 'g',
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
			rev:      5,
		}},
		{1, 2, command{
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
			globalPrefix: 'g',
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
			rev:      4,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doGlob ", idx)
		t.Run(tName, func(t *testing.T) {
			input, _ := newTerm(os.Stdin, os.Stdout, "", false)
			controlBuffer := getTestActionBuffer()
			doGlob(controlBuffer, test.line1, test.line2, test.cmd, input, &localFS{}, &cache{})

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
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

	for idx, test := range tests {
		tName := fmt.Sprint("doSetMarkLine ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			doSetMarkLine(controlBuffer, test.l1, test.l2, test.argument)

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
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

	for idx, test := range tests {
		tName := fmt.Sprint("doSetMarkLine2 ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestMarkedActionBuffer()
			doSetMarkLine(controlBuffer, test.l1, test.l2, test.argument)

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}

func Test_doJoinLines(t *testing.T) {
	tests := []struct {
		l1, l2   int
		expected buffer
	}{
		{2, 4, &memoryBuf{
			curline:  2,
			lastline: 3,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...++3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...++4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
			},
			filename: "filename",
			dirty:    true,
			rev:      9,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doJoinLines ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			doJoinLines(controlBuffer, test.l1, test.l2, `++`)

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}

func Test_doBreakLines_n1(t *testing.T) {
	tests := []struct {
		l1, l2   int
		expected buffer
	}{
		{2, 4, &memoryBuf{
			curline:  4,
			lastline: 8,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut porta mi,`},
				{txt: ` eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor.`},
				{txt: ` Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna,`},
				{txt: ` congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`}, // this is flotsam
			},
			filename: "filename",
			dirty:    true,
			rev:      36,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doBreakLines_n1 ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			doBreakLines(controlBuffer, test.l1, test.l2, replace{`[,.]`, "", `1`})

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}

func Test_doBreakLines_g(t *testing.T) {
	tests := []struct {
		l1, l2   int
		expected buffer
	}{
		{2, 4, &memoryBuf{
			curline:  4,
			lastline: 22,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut porta mi,`},
				{txt: ` eu ornare orci.`},
				{txt: ` Etiam sed vehicula orci.`},
				{txt: ` .`},
				{txt: `.`},
				{txt: `.`},
				{txt: ``},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor.`},
				{txt: ` Donec pulvinar leo urna,`},
				{txt: ` id .`},
				{txt: `.`},
				{txt: `.`},
				{txt: ``},
				{txt: `4 Nullam lacus magna,`},
				{txt: ` congue aliquam luctus ac,`},
				{txt: ` faucibus vel purus.`},
				{txt: ` Integer .`},
				{txt: `.`},
				{txt: `.`},
				{txt: ``},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`}, // this is flotsam
			},
			filename: "filename",
			dirty:    true,
			rev:      92,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doBreakLines_g ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			doBreakLines(controlBuffer, test.l1, test.l2, replace{`[,.]`, "", `-1`})

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}

func Test_doPrintAddress(t *testing.T) {
	tests := []struct {
		l2       int
		expected string
	}{
		{5, "5"},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doPrintAddress ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			s, err := doPrintAddress(controlBuffer, test.l2)
			if err != nil {
				t.Error(err)
			}

			if s != test.expected {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}

func Test_setFilename(t *testing.T) {
	tests := []struct {
		given    string
		expected buffer
	}{
		{"new file name", &memoryBuf{
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
			filename: "new file name",
			dirty:    false,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doSetFilename ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			buf := new(localFile)

			doSetFilename(controlBuffer, &localFS{buf}, `new file name`)

			if test.given != filepath.Base(test.expected.getFilename()) {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}

func Test_doMirrorLines(t *testing.T) {
	tests := []struct {
		l1, l2   int
		expected buffer
	}{
		{2, 4, &memoryBuf{
			curline:  4, // NOTE: curious
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
			rev:      1,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doMirrorLines ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			doMirrorLines(controlBuffer, test.l1, test.l2)

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}

func Test_doTransliterate(t *testing.T) {
	tests := []struct {
		l1, l2   int
		expected buffer
	}{
		{2, 4, &memoryBuf{
			curline:  1,
			lastline: 5,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut port1 mi, 2u orn1r2 orci. Eti1m s2d v2hicul1 orci. ...`},
				{txt: `3 Nunc sc2l2risqu2 urn1 1 2r1t gr1vid1 porttitor. Don2c pulvin1r l2o urn1, id ...`},
				{txt: `4 Null1m l1cus m1gn1, congu2 1liqu1m luctus 1c, f1ucibus v2l purus. Int2g2r ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
			rev:      3,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doTransliterate ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			doTransliterate(controlBuffer, test.l1, test.l2, `ae`, `12`)

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}

func Test_doPrint(t *testing.T) {
	tests := []struct {
		l1, l2   int
		tp       int
		expected string
	}{
		{1, 2, printTypeReg, "1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...\n2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...\n"},
		{1, 2, printTypeNum, "  1	1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...\n  2	2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...\n"},
		{1, 2, printTypeLit, "  1	\"1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...\"\n  2	\"2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...\"\n"},
		{1, 2, printTypeCol, "  1	1 Lorem ipsum dolor sit amet, consectetur adipisc\x1b[7mi\x1b[0mng elit. Morbi sed ante eu ...\n  2	2 Duis ut porta mi, eu ornare orci. Etiam sed veh\x1b[7mi\x1b[0mcula orci. ...\n"},
	}

	// line, err := term.input(":")
	cache := &cache{
		pager:  0,
		column: 50,
	}

	// t.Error(out.String(), line, err)

	for idx, test := range tests {
		tName := fmt.Sprint("doPrint ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			out := bytes.NewBufferString(``)
			in := bytes.NewBufferString(``)
			term, _ := newTerm(in, out, "", false) // _ is an unused destructor

			doPrint(term, controlBuffer, test.l1, test.l2, cache, test.tp)

			if out.String() != test.expected {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}

func Test_doAppend(t *testing.T) {
	tests := []struct {
		l1       int
		expected buffer
	}{
		{1, &memoryBuf{
			curline:  2,
			lastline: 6,
			lines: []bufferLine{
				{txt: ``},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `foobar snafu`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
			rev:      4,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doAppend ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			out := bytes.NewBufferString(``)
			in := bytes.NewBufferString("foobar snafu\n.\n")
			term, _ := newTerm(in, out, "", false) // _ is an unused destructor

			doAppend(term, controlBuffer, test.l1)

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}

func Test_doInsert(t *testing.T) {
	tests := []struct {
		l1       int
		expected buffer
	}{
		{1, &memoryBuf{
			curline:  1,
			lastline: 6,
			lines: []bufferLine{
				{txt: ``},
				{txt: `foobar snafu`},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "filename",
			dirty:    true,
			rev:      4,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doInsert ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			out := bytes.NewBufferString(``)
			in := bytes.NewBufferString("foobar snafu\n.\n")
			term, _ := newTerm(in, out, "", false) // _ is an unused destructor

			doInsert(term, controlBuffer, test.l1)

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}

func Test_doReadFile(t *testing.T) {
	tests := []struct {
		l1       int
		expected buffer
	}{
		{1, &memoryBuf{
			curline:  1,
			lastline: 6,
			lines: []bufferLine{
				{txt: ``},
				{txt: `This is the contents of the memory file.`},
				{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
				{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
				{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
				{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
				{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
			},
			filename: "example/short",
			dirty:    true,
			rev:      4,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doReadFile ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			seed := bytes.NewBufferString(`This is the contents of the memory file.`)

			doReadFile(controlBuffer, 0, &localFS{
				seed: &localFile{*seed},
			}, "example/short")

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}

func Test_doWriteFile(t *testing.T) {
	tests := []struct {
		l1       int
		expected buffer
	}{
		{1, &memoryBuf{
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
			rev:      0,
		}},
	}

	for idx, test := range tests {
		tName := fmt.Sprint("doWriteFile ", idx)
		t.Run(tName, func(t *testing.T) {
			controlBuffer := getTestActionBuffer()
			out := bytes.NewBufferString(``)
			in := bytes.NewBufferString("")
			term, _ := newTerm(in, out, "", false) // _ is an unused destructor
			buf := new(localFile)

			doWriteFile(term, controlBuffer, 1, 1, &localFS{buf}, "filename")

			if buf.Len() != 379 {
				t.Errorf("write buffer failed: %d", buf.Len())
			}

			if controlBuffer.String() != test.expected.String() {
				t.Errorf("\n-got:  %s\n+want: %s", controlBuffer, test.expected)
			}
		})
	}
}
