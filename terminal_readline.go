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
		AutoComplete:           fCompleter{},
		InterruptPrompt:        ".",
		EOFPrompt:              "q",
		DisableAutoSaveHistory: true,
		VimMode:                false,
		DisableBell:            true,
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

	t.rl.Config.AutoComplete = fCompleter{}

	if prompt == "" {
		t.rl.Config.AutoComplete = &readline.TabCompleter{}
	}

	// t.rl.SetConfig(cfg)
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

type fCompleter struct{}

func (a fCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
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

type nullCompleter struct{}

func (a nullCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	return nil, 0
}

// TabCompleter allows tabs in typed input; a little copying is better than a little dependency
type TabCompleter struct{}

func (t TabCompleter) Do([]rune, int) ([][]rune, int) {
	return [][]rune{[]rune("\t")}, 0
}
