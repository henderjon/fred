//go:build readline

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"

	"github.com/chzyer/readline"
)

type readlineTerm struct {
	in      *bufio.Scanner
	rl      *readline.Instance
	out     io.Writer
	history []string
	prompt  string
	isPipe  bool
}

// func newLocalTerm(raw bool, in io.ReadWriter, out io.Writer) *localTerm {
func newTerm(pipe io.Reader, stdout io.Writer, prompt string, isPipe bool) (termio, func()) {

	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 prompt,
		InterruptPrompt:        ".",
		EOFPrompt:              "q",
		DisableAutoSaveHistory: true,
		VimMode:                false,
	})
	if err != nil {
		return nil, nil
	}

	stdin := bufio.NewScanner(pipe)
	return &readlineTerm{
			in:     stdin,
			rl:     rl,
			out:    stdout,
			prompt: prompt,
			isPipe: isPipe,
		}, func() {
			rl.Close()
		}
}

func (t *readlineTerm) getPrompt() string {
	return t.prompt
}

func (t *readlineTerm) Write(b []byte) (int, error) {
	b = bytes.TrimRightFunc(b, unicode.IsSpace)
	return fmt.Fprintln(t.out, string(b))
}

func (t *readlineTerm) input(prompt string) (string, error) {
	if t.isPipe {
		prompt = ""
	}

	t.rl.SetPrompt(prompt)
	line, err := t.rl.Readline()

	if prompt != "" {
		t.rl.SaveHistory(line)
	}

	return line, err
}

func (t *readlineTerm) prtHistory(s string) error {
	return errors.New("command not implemented for this terminal")
}