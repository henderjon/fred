package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"

	"golang.org/x/term"
)

func main() {
	// switch stdin into 'raw' mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for i := 0; i < 10; i++ {
		b := make([]byte, 4)
		_, err := os.Stdin.Read(b)
		if err != nil {
			fmt.Println(err)
			return
		}

		// https://www2.ccs.neu.edu/research/gpc/VonaUtils/vona/terminal/vtansi.htm
		fmt.Fprint(os.Stdout, "\0337")
		fmt.Fprintf(
			os.Stderr,
			"the char %q or hex(% x) dec(%v) type(%s) was hit\n",
			string(b),
			b,
			b,
			getRange(b2r(b)),
		)
		fmt.Fprint(os.Stdout, "\0338")
		// fmt.Fprint(os.Stdout, "\033[1B")
	}
}

func b2r(bts []byte) rune {
	b := bytes.TrimRight(bts, "\x00")
	s := string(b)
	return rune(s[0])
}

func getRange(r rune) string {
	funcs := map[string]func(rune) bool{
		"IsControl": unicode.IsControl,
		"IsDigit":   unicode.IsDigit,
		"IsGraphic": unicode.IsGraphic,
		"IsLetter":  unicode.IsLetter,
		"IsLower":   unicode.IsLower,
		"IsMark":    unicode.IsMark,
		"IsNumber":  unicode.IsNumber,
		"IsPrint":   unicode.IsPrint,
		"IsPunct":   unicode.IsPunct,
		"IsSpace":   unicode.IsSpace,
		"IsSymbol":  unicode.IsSymbol,
		"IsTitle":   unicode.IsTitle,
		"IsUpper":   unicode.IsUpper,
	}

	types := make([]string, 0)

	for name, fn := range funcs {
		if fn(r) {
			types = append(types, name)
		}
	}
	return strings.Join(types, ",")
}
