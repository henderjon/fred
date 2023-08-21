package main

type termio interface {
	Fprintf(format string, args ...any) (int, error)
	Fprintln(args ...any) (int, error)
	input(prompt string) (string, error)
}
