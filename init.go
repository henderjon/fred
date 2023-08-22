package main

import (
	"flag"
	"os"
)

type generalParams struct {
	debug    bool
	filename string
	prompt   string
	// gutter   int
	pager        int
	experimental bool
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
	flag.StringVar(&params.general.filename, "file", "", "edit `filename`; the last unnamed arg will be used if not provided")
	flag.StringVar(&params.general.prompt, "prompt", ":", "the string to display at the beginning of each line")
	// flag.IntVar(&params.general.gutter, "gutter", 2, "the space-padded width of the line number gutter")
	flag.IntVar(&params.general.pager, "pager", 0, "the space-padded width of the line number gutter")
	flag.BoolVar(&params.general.experimental, "experimental", false, "fancy terminal [not operational]")
	flag.Parse()

	if params.general.debug {
		os.Exit(0)
	}

	// args := flag.Args()
	// if len(args) > 0 {
	// 	params.general.infile = args[len(args)-1]
	// }

	// if len(params.general.infile) > 0 {
	// 	filename := params.general.infile
	// 	doRead(0, filename)
	// }

	return params
}
