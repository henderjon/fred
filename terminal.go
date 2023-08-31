package main

import "io"

type termio interface {
	io.Writer
	getPrompt() string
	input(prompt string) (string, error)
	prtHistory(string) error
}
