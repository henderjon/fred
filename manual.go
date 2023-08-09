package main

import (
	"bytes"
	"flag"
	"os"
	"text/template"
)

// these vars are built at compile time, DO NOT ALTER
var (
	// Version adds build information
	binName string
	// Version adds build information
	buildVersion string
	// BuildTimestamp adds build information
	buildTimestamp string
	// CompiledBy adds the make/model that was used to compile
	compiledBy string
)

// Tmpl is a basic man page[-ish] looking template
const Tmpl = `
{{define "manual"}}
NAME
  {{.Bin}} - FRiendlier ED; another attempt to recreate ed(1)

SYNOPSIS
  $ {{.Bin}}
  $ {{.Bin}} [-h|help]

DESCRIPTION
  {{.Bin}} is a terminal line-based text editor akin to ed(1). {{.Bin}} shares
  as much feature parity as possible with the ed(1) implementation found on
  macOS which, as is understood by the author, is closer to BSD's ed(1) than the
  GNU version.

  The differences found follow the author's usage patterns. In other words, the
  differences between {{.Bin}} and the original ed(1) are based on how the
  author uses ed(1).

  For detailed documentation on how to use {{.Bin}}/ed(1), use 'man ed' to read
  the ed(1) documentation. This document will attempt at only capturing the
  differences.

WHY
  {{.Bin}}'s goals were:

  - Understand ed(1)
  - Create a re-readable code base not written by a proper genius
  - Use as few external libs as possible
  - Preserve the spirit of ed(1) as fewer and fewer common devs work in
    advanced, genius-level C
  - Avoid global state even at cost

DIFFERENCES WITH THE ORIGINAL ED
  - As of now, {{.Bin}} does not support the interactive elements of 'G' and 'V'
  - In {{.Bin}}, the 'z' (scroll) command is used to set a 'pager' that displays
    n line(s) before and after the current line. This has the added effect of
	moving the current line forward n lines. This is applied to all un-bound
	print commands.
  - In {{.Bin}}, 'H' isn't necessary as all errors are displayed as hiding them
    doesn't seem all the necessary.
  - Some of the default line/ranges have been tweaked to reflect how the author
    uses ed(1).
  - Showing line numbers is default.
  - Marking a line ('k') is omitted for now. Users can mark lines with ' and
    retrieve them with ". Marking lines takes a single character as a label.

NOTES
  - During development, it came to light that using Rob Pike's lexer/parser
    would not only yield non-global readable code, but it would also lend the
	  code base to being testable.
  - There are times when verbosity was chosen over brevity for the sake of
    clarity. Certainly, there are places where this code takes the long way
    around ("brute force" so to speak). Perhaps over time this code base will
    evolve.

EXAMPLES
  $ {{.Bin}} -h

OPTIONS
{{.Options}}
VERSION
  version:  {{.Version}}
  compiled: {{.CompiledBy}}
  built:    {{.BuildTimestamp}}

{{end}}
`

// Info represents the infomation used in the default Tmpl string
type Info struct {
	Tmpl           string
	Bin            string
	Version        string
	CompiledBy     string
	BuildTimestamp string
	Options        string
}

// Usage wraps a set of `Info` and creates a flag.Usage func
func Usage(info Info) func() {
	if len(info.Tmpl) == 0 {
		info.Tmpl = Tmpl
	}

	t := template.Must(template.New("manual").Parse(info.Tmpl))

	return func() {
		var def bytes.Buffer
		flag.CommandLine.SetOutput(&def)
		flag.PrintDefaults()

		info.Options = def.String()
		t.Execute(os.Stdout, info)
	}
}
