package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envlens/internal/audit"
	"github.com/yourorg/envlens/internal/diff"
	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/report"
)

func main() {
	var (
		format    = flag.String("format", "text", "Output format: text or json")
		auditMode = flag.Bool("audit", false, "Run audit on diff results")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envlens [flags] <base-file> <target-file>\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	baseFile, targetFile := args[0], args[1]

	baseEnv, err := parser.ParseFile(baseFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading base file %q: %v\n", baseFile, err)
		os.Exit(1)
	}

	targetEnv, err := parser.ParseFile(targetFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading target file %q: %v\n", targetFile, err)
		os.Exit(1)
	}

	entries := diff.Compare(baseEnv, targetEnv)

	if *auditMode {
		findings := audit.Audit(entries)
		switch *format {
		case "json":
			fmt.Print(report.AuditJSONReport(findings))
		default:
			fmt.Print(report.AuditTextReport(findings))
		}
		return
	}

	switch *format {
	case "json":
		fmt.Print(report.JSONReport(entries))
	default:
		fmt.Print(report.TextReport(entries))
	}
}
