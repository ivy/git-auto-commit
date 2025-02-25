package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"

	// Assume this package registers flags for --provider, --model, --openai-key, etc.
	"github.com/ivy/git-auto-commit/config"
)

// Program metadata
const (
	ProgramName = "git auto-pr"
	RepoURL     = "https://github.com/ivy/git-auto-commit"
)

// Version is set via build flags (e.g. `-ldflags="-X main.Version=1.0.0"`).
// Defaults to "dev" if not set.
var Version = "dev"

// CLIFlags holds local CLI-only flags that are *not* in config.Config.
type CLIFlags struct {
	Verbose bool
	Yes     bool
	Message string
}

func main() {
	// 1. Initialize the config package so it can register pflags
	config.Init()

	// 2. Local flags for this CLI
	var (
		cli     CLIFlags
		showVer bool
	)

	// 3. Customize pflag usage to display program help
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
  # Add a custom message and open the PR after creation:
  %s -m "My custom message" -- --open

Options:
`,
			ProgramName, Version, RepoURL, os.Args[0], os.Args[0],
		)
		pflag.PrintDefaults()
	}

	// 4. Register local flags that aren't part of config.Config
	pflag.BoolVarP(
		&showVer, "version", "V", false,
		"Print the version of this tool and exit.",
	)
	pflag.BoolVarP(
		&cli.Verbose, "verbose", "v", false,
		"Opens your $EDITOR with a suggested PR description.",
	)
	pflag.BoolVarP(
		&cli.Yes, "yes", "y", false,
		"Creates your PR with the suggested message without prompting.",
	)
	pflag.StringVarP(
		&cli.Message, "message", "m", "",
		"Adds extra context for the LLM (why the change was made).",
	)

	// 5. Parse the flags *once*.
	pflag.Parse()

	// 6. Any leftover arguments become PR arguments for `gh`.
	prArgs := pflag.Args()

	if showVer {
		fmt.Printf("%s %s\n", ProgramName, Version)
		os.Exit(0)
	}

	// 7. Load our layered configuration (from defaults, git config, environment, pflags).
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 8. For demonstration, print out the loaded config fields
	fmt.Fprintf(os.Stderr,
		"Using config:\n  Provider=%s\n  Model=%s\n  OpenAIKey=%q\n\n",
		cfg.Provider, cfg.Model, cfg.OpenAIAPIKey,
	)
	fmt.Fprintf(os.Stderr,
		"Local flags:\n  Verbose=%v\n  Yes=%v\n  Message=%q\n  PRArgs=%v\n",
		cli.Verbose, cli.Yes, cli.Message, prArgs,
	)

	// Actual PR creation logic would go here...
	fmt.Fprintln(os.Stderr, "Error: 'git auto-pr' is not implemented yet.")
	os.Exit(1)
}
