package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/henderjon/shutdown"
)

type generalParams struct {
	debug    bool
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
	flag.BoolVar(&params.general.debug, "debug", false, "output diagnostic information; currently does nothing")
	flag.StringVar(&params.general.filename, "file", "", "load `filename`; but also assumes the last non-flag arg is `filename`")
	flag.StringVar(&params.general.prompt, "prompt", ":", "the string to display at the beginning of each line")
	flag.IntVar(&params.general.pager, "pager", 0, "the number of contextual lines to display before and after the current line when printing")
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
	inout, destructor := newTerm(os.Stdin, os.Stdout, opts.general.prompt, isPipe(os.Stdin))

	shd := shutdown.New(func() {
		b.destructor() // clean up our tmp file
		destructor()   // close readline
	})
	// defer shd.Destructor()

	// HUP signals
	// signal.Notify(sysSigChan, syscall.SIGINT)
	// signal.Notify(sysSigChan, syscall.SIGTERM)
	// signal.Notify(sysSigChan, syscall.SIGHUP)

	if len(opts.general.filename) > 0 {
		numbts, err := doReadFile(b, b.getCurline(), osFS{}, opts.general.filename)
		if err != nil {
			fmt.Fprintln(inout, err.Error())
		} else {
			fmt.Fprintln(inout, numbts)
			b.setDirty(false) // loading the file on init isn't *actually* dirty
		}
	}

	return shd, inout
}
