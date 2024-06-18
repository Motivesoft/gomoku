package main

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

var version string = "0.0.0 (undefined)"

func main() {
	// Configure the flags
	showVersion := flag.BoolP("version", "v", false, "show version")
	showHelp := flag.BoolP("help", "h", false, "show help")

	flag.Usage = func() {
		flag.PrintDefaults()
	}

	flag.Parse()

	if showVersion != nil && *showVersion {
		fmt.Printf("v%s\n", version)
	} else if showHelp != nil && *showHelp {
		flag.Usage()
	} else {
		// Assume we want to build
	}
}
