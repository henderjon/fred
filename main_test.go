//go:build memory

package main

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_print_command(t *testing.T) {
	tests := []struct {
		cmd      *command
		expected string
	}{
		{&command{
			addrStart: "1",
			addrEnd:   "2",
			action:    'p',
		}, "1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...\n2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...\n"},
	}

	cache := &cache{}

	for i, test := range tests {
		controlBuffer := getTestActionBuffer()
		out := bytes.NewBufferString(``)
		in := bytes.NewBufferString(``)
		term, _ := newTerm(in, out, "", false) // _ is an unused destructor

		doCmd(*test.cmd, controlBuffer, term, &localFS{}, cache)

		// t.Error(controlBuffer.String(), test.expected.String())

		if diff := cmp.Diff(out.String(), test.expected); diff != "" {
			t.Errorf("idx: %d; -got/+want\n%s", i, diff)
		}
	}
}
func Test_quit_command(t *testing.T) {
	tests := []struct {
		cmd      *command
		expected error
	}{
		{
			&command{addrStart: "1", addrEnd: "2", action: 'q'},
			errQuit,
		},
	}

	cache := &cache{}

	for i, test := range tests {
		controlBuffer := getTestActionBuffer()
		out := bytes.NewBufferString(``)
		in := bytes.NewBufferString(``)
		term, _ := newTerm(in, out, "", false) // _ is an unused destructor

		_, err := doCmd(*test.cmd, controlBuffer, term, &localFS{}, cache)

		if diff := cmp.Diff(err.Error(), test.expected.Error()); diff != "" {
			t.Errorf("idx: %d; -got/+want\n%s", i, diff)
		}
	}
}

func Test_interactive_glob(t *testing.T) {
	tests := []struct {
		cmd      *command
		expected buffer
	}{
		{
			// &command{addrPattern: "ui", action: 's', pattern: "ui", replace: "++", replaceNum: "g"},
			&command{addrPattern: "ui", globalPrefix: 'G'},
			&memoryBuf{
				curline:  2,
				lastline: 5,
				lines: []bufferLine{
					{txt: ``},
					{txt: `1 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi sed ante eu ...`},
					{txt: `2 D++s ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
					{txt: `3 Nunc scelerisque urna a erat gravida porttitor. Donec pulvinar leo urna, id ...`},
					{txt: `4 Nullam lacus magna, congue aliquam luctus ac, faucibus vel purus. Integer ...`},
					{txt: `5 Mauris nunc purus, congue non vehicula eu, blandit sit amet est. ...`},
				},
				filename: "filename",
				dirty:    true,
				rev:      1,
			},
		},
	}

	cache := &cache{}

	for i, test := range tests {
		controlBuffer := getTestActionBuffer()
		out := bytes.NewBufferString(``)
		in := bytes.NewBufferString("s/ui/++/g\nQ\n")
		term, _ := newTerm(in, out, "", false) // _ is an unused destructor

		_, err := doCmd(*test.cmd, controlBuffer, term, &localFS{}, cache)
		if err != nil && err != errStop {
			t.Error(err)
		}

		if diff := cmp.Diff(controlBuffer.String(), test.expected.String()); diff != "" {
			t.Errorf("idx: %d; -got/+want\n%s", i, diff)
		}
	}
}
func Test_interactive_v_glob(t *testing.T) {
	tests := []struct {
		cmd      *command
		expected buffer
	}{
		{
			// &command{addrPattern: "ui", action: 's', pattern: "ui", replace: "++", replaceNum: "g"},
			&command{addrPattern: "ui", globalPrefix: 'V'},
			&memoryBuf{
				curline:  2,
				lastline: 5,
				lines: []bufferLine{
					{txt: ``},
					{txt: `1 Lor!m ipsum dolor sit am!t, cons!ct!tur adipiscing !lit. Morbi s!d ant! !u ...`},
					{txt: `2 Duis ut porta mi, eu ornare orci. Etiam sed vehicula orci. ...`},
					{txt: `3 Nunc sc!l!risqu! urna a !rat gravida porttitor. Don!c pulvinar l!o urna, id ...`},
					{txt: `4 Nullam lacus magna, congu! aliquam luctus ac, faucibus v!l purus. Int!g!r ...`},
					{txt: `5 Mauris nunc purus, congu! non v!hicula !u, blandit sit am!t !st. ...`},
				},
				filename: "filename",
				dirty:    true,
				rev:      4,
			},
		},
	}

	cache := &cache{}

	for i, test := range tests {
		controlBuffer := getTestActionBuffer()
		out := bytes.NewBufferString(``)
		in := bytes.NewBufferString("s/e/!/g\nq\ns/e/!/g\nq\ns/e/!/g\nq\ns/e/!/g\nq\nQ\n")
		term, _ := newTerm(in, out, "", false) // _ is an unused destructor

		_, err := doCmd(*test.cmd, controlBuffer, term, &localFS{}, cache)
		if err != nil && err != errStop {
			t.Error(err)
		}

		if diff := cmp.Diff(controlBuffer.String(), test.expected.String()); diff != "" {
			t.Errorf("idx: %d; -got/+want\n%s", i, diff)
		}
	}
}
func Test_doExternalShell(t *testing.T) {
	tests := []struct {
		cmd            *command
		expectedResult string
		expectedReturn string
	}{
		{
			&command{action: '!', argument: "echo 867-5309"},
			"867-5309\n",
			"!",
		},
	}

	cache := &cache{}

	for i, test := range tests {
		controlBuffer := getTestActionBuffer()
		out := bytes.NewBufferString(``)
		in := bytes.NewBufferString("s/e/!/g\nq\ns/e/!/g\nq\ns/e/!/g\nq\ns/e/!/g\nq\nQ\n")
		term, _ := newTerm(in, out, "", false) // _ is an unused destructor

		str, err := doCmd(*test.cmd, controlBuffer, term, &localFS{}, cache)
		if err != nil {
			t.Error(err)
		}

		if diff := cmp.Diff(out.String(), test.expectedResult); diff != "" {
			t.Errorf("idx: %d; -got/+want\n%s", i, diff)
		}

		if diff := cmp.Diff(str, test.expectedReturn); diff != "" {
			t.Errorf("idx: %d; -got/+want\n%s", i, diff)
		}
	}
}
