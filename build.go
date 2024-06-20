package main

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

func commandBuild(commandArgs []string) error {
	var showHelp bool = false
	var buildAll bool = false

	var goos string
	var goarch string

	buildCommand := flag.NewFlagSet("build", flag.ExitOnError)
	buildCommand.BoolVarP(&showHelp, "help", "h", false, "show help for this command")

	// Disable this until we can offer it?
	buildCommand.BoolVarP(&buildAll, "all", "a", false, "build for all platforms")

	// These must be either both present or both absent
	buildCommand.StringVar(&goos, "goos", "", "Target platform")
	buildCommand.StringVar(&goarch, "goarch", "", "Target architecture")

	buildCommand.Usage = func() {
		_, executable := filepath.Split(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage information for build command:\n")
		fmt.Fprintf(os.Stderr, "  %s build [arguments]\n\n", executable)
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, " --help, -h                  show this help\n")
		fmt.Fprintf(os.Stderr, " --all, -a                   build for all selected platforms\n")
		fmt.Fprintf(os.Stderr, " --goos <platform>           build for specified platform\n")
		fmt.Fprintf(os.Stderr, " --goarch <architecture>     build for specified architecture\n")
	}

	err := buildCommand.Parse(commandArgs)
	if err != nil {
		return err
	}

	if showHelp {
		buildCommand.Usage()
		return nil
	}

	// Validate the options
	if goos == "" {
		if goarch != "" {
			return fmt.Errorf("use of --goarch also requires --goos")
		}
	} else if goarch == "" {
		return fmt.Errorf("use of --goos also requires --goarch")
	} else if buildAll {
		return fmt.Errorf("cannot use --all with --goos and --goarch")
	}

	var moduleName string = "."
	if buildCommand.NArg() > 0 {
		moduleName = buildCommand.Arg(0)
	}

	fmt.Printf("Building '%s'", moduleName)
	if buildAll {
		fmt.Printf(" for all platforms\n")
	} else if goos != "" && goarch != "" {
		fmt.Printf(" for %s/%s\n", goos, goarch)
	} else {
		fmt.Printf(" for current platform\n")
	}

	return nil
}
