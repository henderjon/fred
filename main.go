package main

import (
	"os"

	"github.com/henderjon/logger"
)

var (
	stderr = logger.NewDropLogger(os.Stderr)
)

func main() {
	input := `15s/pattern/substitute/`
	c, err := (&parser{}).run(input)
	if err != nil {
		stderr.Log(err)
	}
	stderr.Log(input, c.String())
}
