package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

type localTerm struct {
	eol  string
	term *term.Terminal
}

// func newLocalTerm(raw bool, in io.ReadWriter, out io.Writer) *localTerm {
func newLocalTerm(in *os.File, out io.Writer) (termio, func()) {
	oldState := makeTerminal(in)

	var stdin io.ReadWriter = in // this will freak out if `in` isn't an io.ReadWriter

	t := term.NewTerminal(stdin, "")
	return &localTerm{
			eol:  "\r\n",
			term: t,
		}, func() {
			term.Restore(int(in.Fd()), oldState)
		}
}

func makeTerminal(in *os.File) *term.State {
	oldState, err := term.MakeRaw(int(in.Fd()))
	if err != nil {
		panic(err)
	}
	return oldState
}

func (t *localTerm) printf(format string, a ...any) (n int, err error) {
	s := fmt.Sprintf(format, a...)
	return t.println(s)
}

func (t *localTerm) println(a ...any) (n int, err error) {
	s := fmt.Sprint(a...)
	s = strings.TrimRight(s, "\r\n")
	return fmt.Fprint(t.term, s, t.eol)
}

func (t *localTerm) input(prompt string) (string, error) {
	t.term.SetPrompt(prompt)
	ln, err := t.term.ReadLine()
	return string(ln), err
}
