package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func main() {

	fmt.Println("Starting terminal mode.")
	fmt.Println("Enter h for [h]elp.")
	fmt.Println("Enter l for [l]ist of commands.")
	fmt.Println("Enter q for [q]uit.")

	termState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set raw mode on STDIN: %v\n", err)
		return
	}

	n := term.NewTerminal(os.Stdin, ": ")
	n.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		k := string(key)
		if len(k) == 3 {
			if k == "\x1b[A" || // up
				k == "\x1b[C" || // right
				k == "\x1b[B" || //down
				k == "\x1b[D" { //left
				fmt.Print(key)
				// return line, pos, true
			}
		}

		// if unicode.IsControl(key) {
		return line, pos, false
	}

	var ln string
	for i := 0; i < 10; i++ {
		ln, err = n.ReadLine()

		if err != nil {
			break
		}
		fmt.Fprintf(os.Stdout, (ln))
		fmt.Fprintf(os.Stdout, "\033[1S\033[1E")
	}
	term.Restore(int(os.Stdin.Fd()), termState)
}

func handleTabs(s string) string {
	s = strings.ReplaceAll(s, `\\t`, "\x1A")
	s = strings.ReplaceAll(s, `\t`, "\t")
	return strings.ReplaceAll(s, "\x1A", `\t`)
}
