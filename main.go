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

	flag.SetInterspersed(false)
	flag.Parse()

	// Process the command line options or command
	var err error
	if showVersion {
		fmt.Printf("v%s\n", version)
	} else if showHelp || flag.NArg() == 0 {
		flag.Usage()
	} else {
		// Determine and execute the subcommand
		switch flag.Args()[0] {
		case "build":
			err = commandBuild(flag.Args()[1:])

		default:
			err = fmt.Errorf("unknown command: %s", flag.Args()[0])
		}
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func commandBuild(commandArgs []string) error {
	var buildAll bool = false

	var goos string
	var goarch string

	buildCommand := flag.NewFlagSet("build", flag.ExitOnError)
	buildCommand.BoolVarP(&buildAll, "all", "a", false, "build for all platforms")

	buildCommand.StringVar(&goos, "goos", "", "Target platform")
	buildCommand.StringVar(&goarch, "goarch", "", "Target architecture")

	buildCommand.Usage = func() {
		_, executable := filepath.Split(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage information for build command:\n")
		fmt.Fprintf(os.Stderr, "  %s build [arguments]\n\n", executable)
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, " --all, -a     build for all platforms\n")
	}

	err := buildCommand.Parse(commandArgs)
	if err != nil {
		return err
	}

	if buildAll {
		fmt.Println("building for all platforms")
	}

	return nil
}
