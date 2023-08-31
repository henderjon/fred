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
  {{.Bin}} - FRiendlier ED; another attempt to recreate ed (ed is the standard
  Unix text editor)

SYNOPSIS
  $ {{.Bin}}
  $ {{.Bin}} [-h|help]

DESCRIPTION
  {{.Bin}} is a terminal line-based text editor akin to ed. {{.Bin}} shares
  as much feature parity as possible with the ed implementation found on
  macOS which, as is understood by the author, is closer to BSD's ed than the
  GNU version.

  Because 'Clear is better than clever' and the author can't read genius level
  C, brute force was employed in an attempt at clarity.

  The differences found follow the author's usage patterns. In other words, the
  differences between {{.Bin}} and the original ed are based on how the
  author uses ed.

  For detailed documentation on how to use {{.Bin}}/ed, use 'man ed' to read
  the ed documentation. While there are differences between ed and fred,
  the usage is consistent. This document will attempt at only capturing the
  differences.

WHY
  {{.Bin}}'s goals were:

  - Understand ed
  - Create a re-readable code base not written by a proper genius
  - Use as few external libs as possible
  - Preserve the spirit of ed as fewer and fewer common devs work in
    advanced, genius-level C
  - Avoid global state even at cost

DIFFERENCES WITH THE ORIGINAL ED
  - {{.Bin}} supports a raw terminal. Meaning arrows are respected when entering
    commands.
  - As of now, {{.Bin}} does not support the interactive elements of 'G' and 'V'
  - In {{.Bin}}, the 'z' (scroll) command is used to set a 'pager' that displays
    n line(s) before and after the current line. This has the added effect of
    moving the current line forward n lines. This is applied to all un-bound
    print commands.
  - In {{.Bin}}, 'H' is used to show command history not verbose errors.
  - Some of the default line/ranges have been tweaked to reflect how the author
    uses ed.
  - Showing line numbers is default.
  - Marking a line ('k') is omitted for now. Users can mark lines with ' and
    retrieve them with ". Marking lines takes a single character as a label.
  - {{.Bin}} doesn't support arithmetic operators with address numbers (e.g +2
    or -7). Instead, {{.Bin}} uses '<' as a look behind operator and '>' as a
    look ahead operator.

NOTES
  - During development, it came to light that using Rob Pike's lexer/parser
    would not only yield non-global readable code, but it would also lend the
    code base to being testable.
  - There are times when verbosity was chosen over brevity for the sake of
    clarity. Certainly, there are places where this code takes the long way
    around ("brute force" so to speak). Perhaps over time this code base will
    evolve.
  - ed's original author was a proper genius and use recursion as their first
    language. {{.Bin}} is the opposite. The long way around was taken as a way
    of making the code base readable even at the expense of duplication and
    verbosity.

COMMANDS
  Most commands default to the current line but some take a range:
    - (.) indicates a single line
    - (,) indicates a range of lines


  - (.)a   [append] Adds the provided input after the current line

  - (,)b/regex/
           [break]  Breaks the given line(s) after each occurrence of regex

  - (,)c   [change] Removes the given line(s) before adding the provided input

  - (,)d   [delete] Removes the given line(s)

  - e file [edit] Clears the current buffer before loading the given file

  - e !exe [edit] Clears the current buffer before loading the output of exe

  - E file [edit] Acts the same as 'e' but without prompting for unsaved changes

  - E !exe [edit] Acts the same as 'e' but without prompting for unsaved changes

  - f name [file] Sets the filename of the current buffer

  - G/regex/
           [global] Search all lines that match regex prompt for commands for
           each line

  - g/regex/cmd
           [global] Search all lines for regex and do cmd

  - Hn     [history] Prints the last n commands

  - h      [help] Shows this document

  - (.)i   [insert] Adds the provided input before the current line.

  - (,)j/sep/
           [join] Joins the given lines using sep

  - (,)k(.)
           [copy] Duplicates the given line(s) to the given destination

  - (,)l   [print] Prints the given line(s) but exposes hidden chars

  - (,)M   [mirror] Reverses the order of the given line(s)

  - (,)m(.)
           [move] Moves the given line(s) to the given destination

  - (,)n   [print] Prints the given line(s) but exposes line numbers

  - (,)p   [print] Prints the given line(s)

  - q      [quit] Prompt for unsaved changes, then Exit

  - Q      [quit] Exit without prompting for unsaved changes

  - (.)r file
           [read] Loads the given file at the given address

  - (.)r !exe
           [read] Loads the output of exe at the given address

  - (,)s/pat/sub/n
           [substitute] Replaces nth occurrence of pat with sub in the given
           line(s). If n is 'g' or '-1' replace all occurrences. If n is absent
           only do the first occurrence

  - (,)S/regex/sub/n
           [substitute] Replaces nth regex match with sub in the given line(s).
           If n is 'g' or '-1' replace all occurrences. If n is absent only do
           the first occurrence

  - (,)t/find/repl/
           [transliterate] Replace each character in find with it corresponding
           char from repl in the given line(s)

  - u      [undo] Restores the previous buffer. Like ed undo is it's own
           inverse

  - V/regex/
           [global] Search all lines that do not match regex prompt for commands
           for each line

  - v/regex/cmd
           [global] Search all lines that do not match regex and do cmd

  - w [filename]
           [write] write the current buffer to the current filename or filename
           if provided. If provided, sets buffer filename.

  - (,)w !exe
           [write] write the given lines to the stdin of exe

  - zn     [pager] Set pager to n. Pager is the number of lines before and after
           to show when printing lines.

  - !exe   [shell] Execute exe in a shell

  - (,)"a  [mark] Get all the lines marked with 'a' [any single character]

  - /regex/
           [search] Find the next line matching regex. Use '//' to repeat

  - (,)'a  [mark] Mark each of the given lines with 'a' [any single character]

  - \regex\
           [search] Find the previous line matching regex. Use '\\' to repeat

  - =      [print] Show only the line number of the current line

  - ^n     [column] Set the column at the nth character

  - (,)|   [column] Print the given lines displaying the column

  -(,)>n   [peek] Print the given lines and the next n lines

  -(,)<n   [peek] Print the given lines and the previous n lines

  - *      [debug] If {{.Bin}} is given -debug, print diagnostic info about the
           current buffer


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

// Info represents the information used in the default Tmpl string
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
