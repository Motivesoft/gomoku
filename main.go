package main

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

var version string = "X.X.X"

func main() {
	// Configure the flags
	var showVersion bool
	var showHelp bool

	flagBoolVarP(&showVersion, "version", "v", false, "show version")
	flagBoolVarP(&showHelp, "help", "h", false, "show help")

	var buildAll bool = false
	// buildCommand := flag.NewFlagSet("build", flag.ExitOnError)
	// flagsetBoolVarP(&buildAll,"all", "a", false, "build for all platforms")

	flag.Usage = func() {
		_, file := filepath.Split(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", file)
		fmt.Fprintf(os.Stderr, "  %s [options] [command] [arguments]\n\n", file)
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, " --version, -v     show version\n")
		fmt.Fprintf(os.Stderr, " --help, -h        show this help\n")
	}

	flag.SetInterspersed(false)
	flag.Parse()
	var err error

	if showVersion {
		fmt.Printf("v%s\n", version)
	} else if showHelp || flag.NArg() == 0 {
		flag.Usage()
	} else {
		// Determine and execute the subcommand
		switch flag.Args()[0] {
		case "build":
			err = build(buildAll)

		default:
			err = fmt.Errorf("unknown command: %s", flag.Args()[0])
		}
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func flagBoolVarP(flagVar *bool, name string, shorthand string, value bool, usage string) {
	flag.BoolVar(flagVar, name, value, usage)
	flag.BoolVar(flagVar, shorthand, value, usage)
}

func flagsetBoolVarP(flagSet *flag.FlagSet, flagVar *bool, name string, shorthand string, value bool, usage string) {
	flagSet.BoolVar(flagVar, name, value, usage)
	flagSet.BoolVar(flagVar, shorthand, value, usage)
}

func build(buildAll bool) error {
	if buildAll {
		fmt.Println("building for all platforms")
	}
	return nil
}

// func executeCommand(args []string) error {
// 	switch args[0] {
// 	case "build":
// 		return executeBuildCommand(args[1:])

// 	default:
// 		return fmt.Errorf("unknown command: %s", args[0])
// 	}
// }

// func executeBuildCommand(args []string) error {
// 	buildCommand := flag.NewFlagSet("build", flag.ExitOnError)

// 	buildAll := buildCommand.BoolP("all", "a", false, "build for all platforms")

// 	buildCommand.Usage()

// 	err := buildCommand.Parse(args)
// 	if err != nil {
// 		return err
// 	}

// 	if *buildAll {
// 		fmt.Println("building for all platforms")
// 	}

// 	return nil
// }
