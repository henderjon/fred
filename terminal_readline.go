//go:build readline

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/chzyer/readline"
)

type readlineTerm struct {
	in       *bufio.Scanner
	rl       *readline.Instance
	complete *fCompleter // this is gross but readline barfs when I try to toggle completion
	out      io.Writer
	history  []string
	prompt   string
	isPipe   bool
}

// func newLocalTerm(raw bool, in io.ReadWriter, out io.Writer) *localTerm {
func newTerm(pipe io.Reader, stdout io.Writer, prompt string, isPipe bool) (termio, func()) {
	completer := &fCompleter{true}
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 prompt,
		AutoComplete:           completer,
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
			in:       stdin,
			rl:       rl,
			complete: completer,
			out:      stdout,
			prompt:   prompt,
			isPipe:   isPipe,
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

	t.complete.toggle(true)
	if prompt == "" {
		t.complete.toggle(false) // when we're in commands that take input, turn off file completion
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

type fCompleter struct {
	act bool
}

func (a *fCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	if !a.act {
		return nil, 0
	}
	_, after, _ := strings.Cut(string(line), " ")

	names := make([][]rune, 0)
	files, _ := os.ReadDir("./")
	for _, f := range files {
		fname := f.Name()
		if strings.HasPrefix(fname, after) {
			names = append(names, []rune(fname[len(after):]))
		}
	}
	return names, len(after)
}

func (a *fCompleter) toggle(b bool) {
	a.act = b
}
