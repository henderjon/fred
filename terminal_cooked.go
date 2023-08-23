package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type classicTerm struct {
	in      *bufio.Scanner
	out     io.Writer
	history []string
}

// func newLocalTerm(raw bool, in io.ReadWriter, out io.Writer) *localTerm {
func newClassicTerm(in io.Reader, out io.Writer) (termio, func()) {
	stdin := bufio.NewScanner(in)
	return &classicTerm{
		in:  stdin,
		out: out,
	}, func() {}
}

func (t *classicTerm) printf(format string, a ...any) (n int, err error) {
	s := fmt.Sprintf(format, a...)
	return t.println(s)
}

func (t *classicTerm) println(a ...any) (n int, err error) {
	return fmt.Fprintln(t.out, a...)
}

func (t *classicTerm) input(prompt string) (string, error) {
	fmt.Fprint(t.out, prompt)
	t.in.Scan()
	if t.in.Err() == nil {
		t.history = append(t.history, t.in.Text())
	}
	return t.in.Text(), t.in.Err()
}

func (t *classicTerm) prtHistory(s string) error {
	var err error

	num := 5
	if len(s) > 0 {
		num, err = strconv.Atoi(s)
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
			t.printf("%2d: %s", idx, ln)
		}
	}

	return nil
}
