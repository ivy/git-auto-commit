package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

// Program metadata
const (
	ProgramName = "git auto-pr"
	RepoURL     = "https://github.com/ivy/git-auto-commit"
)

// Version is set via build flags (e.g., `-ldflags="-X main.Version=1.0.0"`).
// Defaults to "dev" if not set.
var Version = "dev"

// Config holds all user-configurable options.
type Config struct {
	Verbose  bool
	Yes      bool
	Message  string
	Model    string
	Provider string
	PRArgs   []string // Extra args to be passed to `gh` (GitHub CLI).
}

func main() {
	var (
		cfg     Config
		showVer bool
	)

	// Override the usage function to include custom help text.
	pflag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`%s %s
An AI-powered Git tool that generates pull request descriptions using LLMs.

Repository:  %s
Author:      Ivy Evans <ivy@ivyevans.net>
License:     ISC License

Usage:
  %s [options] [-- <extra GitHub CLI args>]

Examples:
  # Add a custom message and open the PR immediately after creation.
  %s -m "My message" -- --open

Options:
`,
			ProgramName, Version, RepoURL, ProgramName, ProgramName,
		)
		pflag.PrintDefaults()
	}

	pflag.BoolVarP(
		&showVer,
		"version", "V",
		false,
		"Print the version of this tool and exit.",
	)

	pflag.BoolVarP(
		&cfg.Verbose,
		"verbose", "v",
		true,
		"Opens your $EDITOR (or falls back to nano/vi) with a suggested PR message.",
	)

	pflag.BoolVarP(
		&cfg.Yes,
		"yes", "y",
		false,
		"Creates your PR with the suggested message without prompting.",
	)

	pflag.StringVarP(
		&cfg.Message,
		"message", "m",
		"",
		"Adds extra context to the LLM (useful for explaining why the change was made).",
	)

	pflag.StringVarP(
		&cfg.Model,
		"model", "M",
		"",
		"Overrides the default model used for PR message generation.",
	)

	pflag.StringVarP(
		&cfg.Provider,
		"provider", "p",
		"",
		"Overrides the default LLM provider.",
	)

	pflag.Parse()

	if showVer {
		fmt.Printf("%s %s\n", ProgramName, Version)
		os.Exit(0)
	}

	// Capture any leftover arguments as PR arguments to be passed to `gh`.
	cfg.PRArgs = pflag.Args()

	// Since we are not implementing the actual PR creation logic here,
	// just print a message and exit for demonstration.
	fmt.Fprintln(os.Stderr, "Error: 'git auto-pr' is not implemented yet.")
	os.Exit(1)
}
