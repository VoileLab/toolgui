// Command toolgui-desktop builds and runs a ToolGUI desktop app using the
// frontend already shipped in the toolgui-wails module.
//
//	go get -tool github.com/voilelab/toolgui/toolgui-wails/cmd/toolgui-desktop
//	go tool toolgui-desktop build ./cmd/myapp
//	go tool toolgui-desktop run ./cmd/myapp
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const usage = `toolgui-desktop builds and runs a ToolGUI desktop app.

Usage:

	toolgui-desktop build [-o file] [-tags tags] [-ldflags flags] [-webkit auto|4.0|4.1] [-console] [package]
	toolgui-desktop run [-tags tags] [-ldflags flags] [-webkit auto|4.0|4.1] [-console] [package] [-- app arguments]

The package defaults to the current directory. Build writes build/bin/<package>
(with .exe on Windows). Run builds a temporary executable and removes it on exit.
Flags must precede the package. -tags adds to desktop,production and
GOFLAGS tags. -ldflags overrides GOFLAGS ldflags. -webkit selects WebKitGTK on
Linux; auto prefers 4.1, then 4.0. -console keeps the console on Windows.

Builds target the host OS and architecture. macOS and Linux need cgo and their
webview development dependencies. Tagged releases include the frontend; no
Node, Yarn, wails.json or Wails CLI is needed. Source checkouts need real assets
from task asset_wails first. Output is a native executable; app bundles,
installers, icons, signing and notarization are outside this command's scope.
`

type buildOptions struct {
	pkg     string
	out     string
	tags    string
	ldflags string
	webkit  string
	console bool
	appArgs []string
}

func main() {
	os.Exit(command(os.Args[1:], os.Stdout, os.Stderr))
}

func command(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	if args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(stdout, usage)
		return 0
	}
	if args[0] != "build" && args[0] != "run" {
		fmt.Fprintf(stderr, "unknown command %q\n\n%s", args[0], usage)
		return 2
	}
	opts, err := parseOptions(args[0], args[1:], stderr)
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}
	if err != nil {
		fmt.Fprintf(stderr, "toolgui-desktop: %v\n", err)
		return 2
	}
	if err = buildOrRun(args[0], opts, stdout, stderr); err != nil {
		fmt.Fprintf(stderr, "toolgui-desktop: %v\n", err)
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() > 0 {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}

func parseOptions(mode string, args []string, stderr io.Writer) (buildOptions, error) {
	opts := buildOptions{pkg: "."}
	flags := flag.NewFlagSet(mode, flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() { fmt.Fprint(stderr, usage) }
	if mode == "build" {
		flags.StringVar(&opts.out, "o", "", "output executable")
	}
	flags.StringVar(&opts.tags, "tags", "", "additional build tags")
	flags.StringVar(&opts.ldflags, "ldflags", "", "Go linker flags")
	flags.StringVar(&opts.webkit, "webkit", "auto", "Linux WebKitGTK version: auto, 4.0 or 4.1")
	flags.BoolVar(&opts.console, "console", false, "keep the Windows console")
	// Keep application arguments out of flag parsing, including arguments that
	// start with a dash. Build never accepts application arguments.
	for i, arg := range args {
		if arg == "--" {
			if mode != "run" {
				return opts, errors.New("application arguments are only supported by run")
			}
			opts.appArgs = args[i+1:]
			args = args[:i]
			break
		}
	}
	if err := flags.Parse(args); err != nil {
		return opts, err
	}
	if flags.NArg() > 1 {
		return opts, errors.New("expected one package; put flags before the package and run arguments after --")
	}
	if flags.NArg() == 1 {
		opts.pkg = flags.Arg(0)
	}
	if strings.HasPrefix(opts.pkg, "-") || strings.HasSuffix(opts.pkg, ".go") {
		return opts, errors.New("expected a Go package directory or import path, not a flag or .go file")
	}
	if opts.webkit != "auto" && opts.webkit != "4.0" && opts.webkit != "4.1" {
		return opts, errors.New("-webkit must be auto, 4.0 or 4.1")
	}
	return opts, nil
}
