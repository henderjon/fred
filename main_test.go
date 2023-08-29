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
		term, _ := newTerm(in, out) // _ is an unused destructor

		doCmd(*test.cmd, controlBuffer, term, cache)

		// t.Error(controlBuffer.String(), test.expected.String())

		if diff := cmp.Diff(out.String(), test.expected); diff != "" {
			t.Errorf("idx: %d; -got/+want\n%s", i, diff)
		}
	}
}
