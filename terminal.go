package main

import (
	"io"
	"os"
)

type termio interface {
	io.Writer
	getPrompt() string
	input(prompt string) (string, error)
	prtHistory(string) error
}

func isPipe(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) == 0
	// bytes, _ := ioutil.ReadAll(os.Stdin)
	// }
}
