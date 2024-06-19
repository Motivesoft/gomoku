# gomoku - a simplified make system for Go

## Aims
* Build/Install for current platform
* Build for other platforms (GOOS, GOARCH)
* Clean
* Test
* Other build commands? lint, tidy, etc
* Obtain/insert git tag 
* Detect and include icon and/or winres folder
* Configuration file for application name, dependencies, linker switches etc...
  * Consider platforms as config option with `build all`
* Uninstall?

### Ambitions
* Publish to GitHub/other
* Make sure this works for modules as well as applications
* VSCode plugin

### Ideas
* Read go.mod/go.sum for information to avoid user entry
  * Use file in same general format for our own configuration?
* Consider accepting multiple module names in some sort of structured build scenario?
  * Need an example to work with
* Each command in its own file?

## Rules
* Expects a module name, but will accept a `.`
  * With `.` (which can also be a default), read `go.mod` for the module name
* 

## Specifics - flesh this out
* Command-based with targets, e.g.:
  * `gomoku clean`
  * `gomoku build --goos=windows --goarch=amd64`
  * `gomoku build all`

### Internals
* Written in Go
* Use `pflag` for command line handling

### Dependencies - add links or whatever works best for this
* `github.com/spf13/pflag`
* go-github?
* gogit?

#### Obtaining dependencies
```shell
go get github.com/spf13/pflag
```

## References - add links
* mage
* goreleaser
* goyek
* make