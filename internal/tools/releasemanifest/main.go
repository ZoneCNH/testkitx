package main

import (
	"errors"
	"flag"
	"io"
	"os"

	"github.com/ZoneCNH/testkitx/internal/cliutil"
)

func main() {
	os.Exit(runCLI(os.Args[0], os.Args[1:], os.Stdout, os.Stderr))
}

func runCLI(name string, args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(stderr)
	out := flags.String("out", "release/manifest/latest.json", "release manifest output path")
	verify := flags.String("verify", "", "verify an existing release manifest instead of generating one")
	requirePassed := flags.Bool("require-passed", false, "require all release checks to be passed during verification")
	requireClean := flags.Bool("require-clean", false, "require a clean git tree during verification")
	expectVersion := flags.String("expect-version", "", "require the manifest version to match this release version during verification")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	if *verify != "" {
		if err := verifyManifest(*verify, *requirePassed, *requireClean, *expectVersion); err != nil {
			return cliutil.PrintCLIError(stderr, err)
		}
		return cliutil.PrintCLIStatus(stdout, "release evidence verified: %s\n", *verify)
	}

	manifest, err := buildManifest()
	if err != nil {
		return cliutil.PrintCLIError(stderr, err)
	}
	if err := writeManifest(*out, manifest); err != nil {
		return cliutil.PrintCLIError(stderr, err)
	}
	return cliutil.PrintCLIStatus(stdout, "generated %s and %s\n", *out, manifestChecksumPath(*out))
}
