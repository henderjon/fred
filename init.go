package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/henderjon/shutdown"
)

type generalParams struct {
	debug    bool
	version  bool
	filename string
	prompt   string
	pager    int
}

type allParams struct {
	general generalParams
}

func getParams() allParams {
	bi := getBuildInfo()
	flag.Usage = Usage(Info{
		Bin:            bi.getBinName(),
		Version:        bi.getBuildVersion(),
		CompiledBy:     bi.getCompiledBy(),
		BuildTimestamp: bi.getBuildTimestamp(),
		QuickHelp:      quickHelp(),
	})

	params := allParams{}
	flag.BoolVar(&params.general.debug, "debug", false, "output diagnostic information; currently does nothing")
	flag.BoolVar(&params.general.version, "version", false, "output version/build information")
	flag.BoolVar(&params.general.version, "v", false, "output version/build information")
	flag.StringVar(&params.general.filename, "file", "", "load `filename`; but also assumes the last non-flag arg is `filename`")
	flag.StringVar(&params.general.prompt, "prompt", ":", "the string to display at the beginning of each line")
	flag.IntVar(&params.general.pager, "pager", 0, "the number of contextual lines to display before and after the current line when printing")
	flag.Parse()

	if params.general.version {
		fmt.Fprint(os.Stdout, getBuildInfo().String())
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) > 0 {
		params.general.filename = args[len(args)-1]
	}

	return params
}

func bootstrap(b buffer, c *cache, opts allParams) (*shutdown.Shutdown, termio) {
	inout, destructor := newTerm(os.Stdin, os.Stdout, opts.general.prompt, isPipe(os.Stdin))

	shd := shutdown.New(nil, []syscall.Signal{
		syscall.SIGHUP,
	})

	shd.SetDestructor(func() {
		if shd.IsDown() { // must create shd before defining this check
			c.getCurrBuffer().destructor() // clean up our tmp file
		}
		destructor() // close readline
	})

	if len(opts.general.filename) > 0 {
		rdr, err := (&osFS{}).FileReader(opts.general.filename)
		if err != nil {
			fmt.Fprintln(inout, err.Error())
		}
		n, err := io.Copy(b, rdr)
		if err != nil {
			fmt.Fprintln(inout, err.Error())
		}

		fmt.Fprintln(inout, n)
		b.setFilename(opts.general.filename)
		b.setDirty(false) // loading the file on init isn't *actually* dirty
	}

	return shd, inout
}
