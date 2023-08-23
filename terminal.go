package main

type termio interface {
	printf(format string, args ...any) (int, error)
	println(args ...any) (int, error)
	input(prompt string) (string, error)
	prtHistory(string) error
}
