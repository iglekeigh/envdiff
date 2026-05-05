package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/envfile"
	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/redact"
)

func main() {
	var (
		redactFlag  = flag.Bool("redact", false, "Redact sensitive values in output")
		formatFlag  = flag.String("format", "text", "Output format: text or json")
		filterFlag  = flag.String("filter", "", "Filter results by status: missing, extra, changed, match")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] <base.env> <other.env>\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	baseFile, otherFile := args[0], args[1]

	baseEnv, err := envfile.Parse(baseFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading base file %q: %v\n", baseFile, err)
		os.Exit(1)
	}

	otherEnv, err := envfile.Parse(otherFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading other file %q: %v\n", otherFile, err)
		os.Exit(1)
	}

	var redactor *redact.Redactor
	if *redactFlag {
		redactor = redact.New()
		baseEnv = redactor.Apply(baseEnv)
		otherEnv = redactor.Apply(otherEnv)
	}

	result := diff.Compare(baseEnv, otherEnv)

	if *filterFlag != "" {
		result = result.Filter(*filterFlag)
	}

	formatter := output.New(*formatFlag)
	if err := formatter.Write(os.Stdout, result); err != nil {
		fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
		os.Exit(1)
	}

	if result.HasDifferences() {
		os.Exit(2)
	}
}
