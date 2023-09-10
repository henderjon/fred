package main

import (
	"fmt"
	"runtime/debug"
	"strings"
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

type buildInfo struct {
	BinName        string
	BuildVersion   string
	BuildTimestamp string
	CompiledBy     string
	BuildTags      string
	LDFlags        string
	ModPath        string
	Rev            string
	RevTimestamp   string
}

func getBuildInfo() buildInfo {
	if binName == "" {
		binName = `fred`
	}

	bi := buildInfo{
		BinName:        binName,
		BuildVersion:   buildVersion,
		BuildTimestamp: buildTimestamp,
		CompiledBy:     compiledBy,
	}

	if bld, ok := debug.ReadBuildInfo(); ok {

		if bi.CompiledBy == "" {
			bi.CompiledBy = bld.GoVersion
		}

		bi.ModPath = bld.Path

		for _, setting := range bld.Settings {
			switch true {
			case setting.Key == "vcs.revision":
				bi.Rev = setting.Value
			case setting.Key == "vcs.time":
				bi.RevTimestamp = setting.Value
			case setting.Key == "-tags":
				bi.BuildTags = setting.Value
			case setting.Key == "-ldflags":
				bi.LDFlags = setting.Value
			}
		}
	}

	return bi
}

func (bi buildInfo) getBinName() string {
	return bi.BinName
}

func (bi buildInfo) getBuildVersion() string {
	return bi.BuildVersion
}

func (bi buildInfo) getBuildTimestamp() string {
	return bi.BuildTimestamp
}

func (bi buildInfo) getCompiledBy() string {
	return bi.CompiledBy
}

func (bi buildInfo) String() string {
	var rtn strings.Builder

	fmt.Fprintf(&rtn, "bin (mod): %s (%s)\n", bi.BinName, bi.ModPath)
	fmt.Fprintf(&rtn, "build version: %s\n", bi.BuildVersion)
	fmt.Fprintf(&rtn, "build timestamp: %s\n", bi.BuildTimestamp)
	fmt.Fprintf(&rtn, "compiled by: %s\n", bi.CompiledBy)
	fmt.Fprintf(&rtn, "build tags: %s\n", bi.BuildTags)

	flags := ""
	if len(bi.LDFlags) > 0 {
		flags = strings.ReplaceAll(bi.LDFlags, "-X ", "\n  -X ")
	}

	fmt.Fprintf(&rtn, "ldflags: %s\n", flags)

	return rtn.String()
}
