package main

import (
	"os"

	"github.com/henderjon/logger"
)

var (
	stderr = logger.NewDropLogger(os.Stderr)
)

func main() {
	// stderr.Log("here")
	c, err := (&parser{}).run("15b")
	if err != nil {
		stderr.Log(err)
	}
	stderr.Log(c.String())
}
