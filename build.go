package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"
)

func commandBuild(commandArgs []string) error {
	var showHelp bool = false
	var buildAll bool = false

	var goos string
	var goarch string

	buildCommand := flag.NewFlagSet("build", flag.ExitOnError)
	buildCommand.BoolVarP(&showHelp, "help", "h", false, "show help for this command")

	// TODO: Disable this until we can offer it?
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

	var module string = "."
	if buildCommand.NArg() > 0 {
		module = buildCommand.Arg(0)
	}

	fmt.Printf("Building '%s'", module)
	if buildAll {
		fmt.Printf(" for all platforms\n")
	} else if goos != "" && goarch != "" {
		fmt.Printf(" for %s/%s\n", goos, goarch)
	} else {
		fmt.Printf(" for current platform\n")
	}

	if module == "." {
		// Read the module from the current directory
		module, err = getGoModule()
		if err != nil {
			return err
		}
	}

	// TODO: Deal with 'all' here by invoking the build function for each platform

	// If goos and goarch are not specified, get their current platform values
	if goos == "" || goarch == "" {
		environment, err := getGoEnvironment()
		if err != nil {
			return err
		}

		goos = environment["GOOS"]
		goarch = environment["GOARCH"]
	}

	return build(module, goos, goarch)
}

func getGoModule() (string, error) {
	output, err := exec.Command("go", "list", "./...").Output()
	if err != nil {
		return "", fmt.Errorf("failed to run 'go list': %s", err)
	}

	return string(output), nil
}

// getGoEnvironment retrieves the Go environment variables and returns them as a map[string]string.
//
// It executes the 'go env' command and parses the output to extract the environment variables.
// The environment variables are stored in a map with the variable name as the key and the value as the value.
// If there is an error executing the 'go env' command, it returns nil and the error.
//
// Returns:
// - map[string]string: A map containing the Go environment variables.
// - error: An error if there was an issue executing the 'go env' command.
func getGoEnvironment() (map[string]string, error) {
	var environment map[string]string = make(map[string]string)

	// We're looking for lines like the following:
	// set GOOS=linux
	regex := *regexp.MustCompile(`set (.*?)=(.*)`)

	// Run the command so we can extract the environment variables
	output, err := exec.Command("go", "env").Output()
	if err == nil {
		res := regex.FindAllStringSubmatch(string(output), -1)
		for i := range res {
			// Extract the environment variable name and value
			environment[res[i][1]] = res[i][2]
		}
	} else {
		return nil, fmt.Errorf("failed to run 'go env': %s", err)
	}

	return environment, nil
}

func build(module, goos string, goarch string) error {
	// Build the executable name
	moduleName := filepath.Base(module)
	executableName := fmt.Sprintf("%s-%s-%s", moduleName, goos, goarch)
	if goos == "windows" {
		executableName = executableName + ".exe"
	}

	// TODO: Make an executable name in a subdirectory?
	//args := fmt.Sprintf("-o %s", filepath.Join(".", "bin", executableName))

	// Build the command line, clause by clause.
	// Start each clause with a space to ensure overall padding
	var args string = "build"
	args += fmt.Sprintf(" -o %s", executableName)
	args += fmt.Sprintf(" %s", module)

	// Now invoke the command with the arguments
	cmd := exec.Command("go", strings.Split(args, " ")...)

	// TODO: Delete this debugging statement
	fmt.Printf("Args:  %s\n", args)

	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()

	return nil
}
