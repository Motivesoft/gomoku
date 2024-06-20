package main

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

var version string = "0.0.0 (development)"

// Prototype of a command function
type commandFunc func([]string) error

// Define a map of command names to their function pointers
var commandMap = map[string]commandFunc{
	"build": commandBuild,
}

func main() {
	// Configure the flags
	var showVersion bool
	var showHelp bool

	flag.BoolVarP(&showVersion, "version", "v", false, "show version")
	flag.BoolVarP(&showHelp, "help", "h", false, "show help")

	flag.Usage = func() {
		_, executable := filepath.Split(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage information for %s:\n", executable)
		fmt.Fprintf(os.Stderr, "  %s [options] [command] [arguments]\n\n", executable)
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, " --version, -v     show version\n")
		fmt.Fprintf(os.Stderr, " --help, -h        show this help\n")
	}

	// Need to set this so that we can parse only the general options that come before the command name
	// Otherwise, the flag parser will get upset by the any additional arguments that should only be passed to the
	// command function
	flag.SetInterspersed(false)
	flag.Parse()

	// Process the command line options or the requested command
	if showVersion {
		fmt.Printf("v%s\n", version)
	} else if showHelp || flag.NArg() == 0 {
		flag.Usage()
	} else {
		// Try and get the command function from the map and invoke it
		commandName := flag.Args()[0]
		commandFunction := commandMap[commandName]

		var err error
		if commandFunction == nil {
			err = fmt.Errorf("unknown command: %s", commandName)
		} else {
			// Invoke the command
			err = commandFunction(flag.Args()[1:])
		}

		// Handle any errors
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
