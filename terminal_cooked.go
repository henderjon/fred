package main

import (
	"bufio"
	"fmt"
	"io"
)

type classicTerm struct {
	in  *bufio.Scanner
	out io.Writer
}

// func newLocalTerm(raw bool, in io.ReadWriter, out io.Writer) *localTerm {
func newClassicTerm(in io.Reader, out io.Writer) (termio, func()) {
	stdin := bufio.NewScanner(in)
	return &classicTerm{
		in:  stdin,
		out: out,
	}, func() {}
}

func (t *classicTerm) Fprintf(format string, a ...any) (n int, err error) {
	s := fmt.Sprintf(format, a...)
	return t.Fprintln(s)
}

func (t *classicTerm) Fprintln(a ...any) (n int, err error) {
	return fmt.Fprintln(t.out, a...)
}

func (t *classicTerm) input(prompt string) (string, error) {
	fmt.Fprint(t.out, prompt)
	t.in.Scan()
	return t.in.Text(), t.in.Err()
}
