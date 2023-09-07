package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

type classicTerm struct {
	in      *bufio.Scanner
	out     io.Writer
	history []string
	prompt  string
	isPipe  bool
}

// func newLocalTerm(raw bool, in io.ReadWriter, out io.Writer) *localTerm {
func newTerm(in io.Reader, out io.Writer, prompt string, isPipe bool) (termio, func()) {
	stdin := bufio.NewScanner(in)
	return &classicTerm{
		in:     stdin,
		out:    out,
		prompt: prompt,
		isPipe: isPipe,
	}, func() {}
}

func (t *classicTerm) getPrompt() string {
	return t.prompt
}

func (t *classicTerm) Write(b []byte) (int, error) {
	b = bytes.TrimRightFunc(b, unicode.IsSpace)
	return fmt.Fprintln(t.out, string(b))
}

func (t *classicTerm) input(prompt string) (string, error) {
	if !t.isPipe {
		fmt.Fprint(t.out, prompt)
	}

	if t.in.Scan() {
		if t.in.Err() == nil && prompt != "" { // skip saving entered text
			t.history = append(t.history, t.in.Text())
		}
		return t.in.Text(), t.in.Err()
	}
	return "", io.EOF
}

func (t *classicTerm) prtHistory(s string) error {
	var err error

	num := 5
	if len(s) > 0 {
		num, err = intval(s)
		if err != nil {
			return err
		}
	}

	x := (len(t.history) - 1) - (num - 1) // get starting index -n from the end
	if x <= 0 {
		x = 0
	}

	for idx, ln := range t.history {
		if idx >= x {
			fmt.Fprintf(t, "%2d: %s", idx, ln)
		}
	}

	return nil
}
