package main

import (
	"flag"
	"os"

	"github.com/henderjon/shutdown"
)

type generalParams struct {
	// debug    bool
	filename string
	prompt   string
	pager    int
}

type allParams struct {
	general generalParams
}

func getParams() allParams {
	flag.Usage = Usage(Info{
		Bin:            binName,
		Version:        buildVersion,
		CompiledBy:     compiledBy,
		BuildTimestamp: buildTimestamp,
	})

	params := allParams{}
	// flag.BoolVar(&params.general.debug, "debug", false, "output the currently set path")
	flag.StringVar(&params.general.filename, "file", "", "load `filename`; but also assumes the last non-flag arg is `filename`")
	flag.StringVar(&params.general.prompt, "prompt", ":", "the string to display at the beginning of each line")
	flag.IntVar(&params.general.pager, "pager", 0, "the space-padded width of the line number gutter")
	flag.Parse()

	// if params.general.debug {
	// 	os.Exit(0)
	// }

	args := flag.Args()
	if len(args) > 0 {
		params.general.filename = args[len(args)-1]
	}

	return params
}

func bootstrap(b buffer, opts allParams) (*shutdown.Shutdown, termio) {
	shd := shutdown.New(func() {
		b.destructor() // clean up our tmp file
	})

	// defer shd.Destructor()

	inout, _ := newTerm(os.Stdin, os.Stdout)

	if len(opts.general.filename) > 0 {
		numbts, err := doReadFile(b, b.getCurline(), opts.general.filename)
		if err != nil {
			inout.println(err.Error())
		} else {
			inout.println(numbts)
			b.setDirty(false) // loading the file on init isn't *actually* dirty
		}
	}

	return shd, inout
}
