package main

import "log"

type LogWriter struct {
	*log.Logger
}

func (lw LogWriter) Write(p []byte) (n int, err error) {
	return lw.Writer().Write(p)
}
