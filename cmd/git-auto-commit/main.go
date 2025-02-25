package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"

	"github.com/ivy/git-auto-commit/config"
)

// Program metadata
const (
	ProgramName = "git auto-commit"
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
	//    for --provider, --model, and --openai-key.
	config.Init()

	// 2. Local flags for this CLI.
	var (
		cli     CLIFlags
		showVer bool
	)

	// Customize pflag usage to display program help.
	pflag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`%s %s
An AI-powered Git tool that generates commit and PR messages using LLMs.

Repository:  %s
Author:      Ivy Evans <ivy@ivyevans.net>
License:     ISC License

Usage:
  %s [options] [-- <extra git commit args>]

Examples:
  # Use GPT-o1, then pass --amend to git commit:
  %s --model=gpt-o1 -- --amend

Options:
`,
			ProgramName, Version, RepoURL, os.Args[0], os.Args[0],
		)
		pflag.PrintDefaults()
	}

	// 3. Register local flags that aren't part of config.Config
	pflag.BoolVarP(
		&showVer, "version", "V", false,
		"Print the version of this tool and exit.",
	)
	pflag.BoolVarP(
		&cli.Verbose, "verbose", "v", false,
		"Opens your $EDITOR with a suggested commit message.",
	)
	pflag.BoolVarP(
		&cli.Yes, "yes", "y", false,
		"Commits changes with the suggested message without prompting.",
	)
	pflag.StringVarP(
		&cli.Message, "message", "m", "",
		"Adds extra context for the LLM (why the change was made).",
	)

	// 4. Parse the pflags *once*.
	pflag.Parse()

	// 5. Any leftover arguments after pflag.Parse() become commit args.
	commitArgs := pflag.Args()

	if showVer {
		fmt.Printf("%s %s\n", ProgramName, Version)
		os.Exit(0)
	}

	// 6. Load our layered config from default, git config, environment, and pflags.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 7. For demonstration, print out the loaded config fields (Provider, Model, etc.)
	//    and your local flags. In a real application, you'd pass them to your logic.
	fmt.Fprintf(os.Stderr,
		"Using config:\n  Provider=%s\n  Model=%s\n  OpenAIKey=%q\n\n",
		cfg.Provider, cfg.Model, cfg.OpenAIAPIKey,
	)
	fmt.Fprintf(os.Stderr,
		"Local flags:\n  Verbose=%v\n  Yes=%v\n  Message=%q\n  CommitArgs=%v\n",
		cli.Verbose, cli.Yes, cli.Message, commitArgs,
	)

	// Actual logic would go here...
	// For now, exit unimplemented:
	fmt.Fprintln(os.Stderr, "Error: 'git auto-commit' is not implemented yet.")
	os.Exit(1)
}
