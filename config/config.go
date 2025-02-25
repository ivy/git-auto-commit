// Package config provides functionality for loading and managing
// user-configurable settings in a layered fashion.
package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/Netflix/go-env"
	"github.com/spf13/pflag"

	"github.com/ivy/git-auto-commit/util/exec"
)

// Config holds user-configurable options. The fields can be set by defaults,
// Git config, environment variables, or command-line flags, in ascending order
// of priority.
type Config struct {
	// Provider denotes the AI provider to use, such as "openai" or "anthropic".
	// By default, this is set to "openai".
	Provider string `env:"GIT_AUTO_COMMIT_PROVIDER"`

	// Model specifies the AI model to use, for example "gpt-4o-mini".
	// By default, this is set to "gpt-4o-mini".
	Model string `env:"GIT_AUTO_COMMIT_MODEL"`

	// OpenAIAPIKey stores the OpenAI token for authentication. This field can
	// only be set via environment variables or pflags, and not from Git config,
	// to avoid checking secrets into Git.
	OpenAIAPIKey string `env:"OPENAI_API_KEY"`
}

// providerFlag, modelFlag, and openAIKeyFlag retain the values passed via the
// corresponding pflags. They are defined here and wired up in Init() so that
// help text is available before Load() is called.
var (
	// providerFlag holds the value of --provider.
	providerFlag *string

	// modelFlag holds the value of --model.
	modelFlag *string

	// openAIKeyFlag holds the value of --openai-key.
	openAIKeyFlag *string
)

// Init registers pflag variables for the Config fields. This function should be
// called before pflag.Parse() to display help text properly.
//
// Example usage in main():
//
//	func main() {
//	    config.Init()
//	    pflag.Parse()
//	    cfg, err := config.Load()
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    // Use cfg...
//	}
func Init() {
	providerFlag = pflag.String("provider", "",
		"Auto-commit provider (overrides env or Git config)")

	modelFlag = pflag.String("model", "",
		"Auto-commit model (overrides env or Git config)")

	openAIKeyFlag = pflag.String("openai-key", "",
		"OpenAI API key (overrides env)")
}

// Load merges configuration from four sources, in ascending priority order:
//  1. Built-in defaults (lowest priority),
//  2. Git config (excluding secret fields),
//  3. Environment variables (via go-env),
//  4. pflag values (highest priority).
//
// It returns a fully populated Config instance, or an error if environment
// unmarshaling fails. Git config is merely logged upon failure, not returned
// as an error.
func Load() (*Config, error) {
	// 1) Built-in defaults.
	cfg := &Config{
		Provider: "openai",
		Model:    "gpt-4o-mini",
	}

	// 2) Git config (non-secret values only).
	getGitConfigValue("auto-commit.provider", &cfg.Provider)
	getGitConfigValue("auto-commit.model", &cfg.Model)
	// We intentionally do not read OpenAIAPIKey from Git config.

	// 3) Environment variables.
	if _, err := env.UnmarshalFromEnviron(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal from environment: %w", err)
	}

	// 4) pflags (highest priority). Override if set to a non-empty value.
	if *providerFlag != "" {
		cfg.Provider = *providerFlag
	}
	if *modelFlag != "" {
		cfg.Model = *modelFlag
	}
	if *openAIKeyFlag != "" {
		cfg.OpenAIAPIKey = *openAIKeyFlag
	}

	return cfg, nil
}

// getGitConfigValue runs `git config --get <key>` to read a single
// configuration value from Git, assigning it to out if successful.
// If any error occurs (like key not found), the error is logged but
// not returned, so that config loading can continue gracefully.
//
// Example:
//
//	getGitConfigValue("auto-commit.provider", &cfg.Provider)
func getGitConfigValue(key string, out *string) {
	raw, err := exec.Command("git", "config", "--get", key).Output()
	if err != nil {
		log.Printf("Error reading git config for %q: %v", key, err)
		return
	}
	trimmed := strings.TrimSpace(string(raw))
	if trimmed != "" {
		*out = trimmed
	}
}
