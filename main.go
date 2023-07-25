package main

import (
	"os"

	"github.com/henderjon/logger"
)

var (
	stderr = logger.NewDropLogger(os.Stderr)
)

func main() {
	// input := `10,15s/pattern/substitute/`
	input := `10,25g/mm/s/and/for/p`
	c, err := (&parser{}).run(input)
	if err != nil {
		stderr.Log(err)
	}

	if c != nil {
		stderr.Log(input, c.String())
	} else {
		stderr.Log(input, "nil command")
	}
}
