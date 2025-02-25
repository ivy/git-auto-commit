package config_test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"

	"github.com/ivy/git-auto-commit/config"
	"github.com/ivy/git-auto-commit/util/exec"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", Serial, func() {
	// We'll hold a copy of the original environment variables
	// and pflag.CommandLine so we can reset between tests.
	var (
		origEnv             []string
		flagSet             *pflag.FlagSet
		logBuf              *bytes.Buffer
		oldFlags            *pflag.FlagSet
		originalExecCommand func(name string, arg ...string) exec.Cmd
	)

	BeforeEach(func() {
		// Save the current environment.
		origEnv = os.Environ()

		// Capture log output to a buffer for easy inspection.
		logBuf = new(bytes.Buffer)
		log.SetOutput(logBuf)

		// Save original execCommand.
		originalExecCommand = exec.Command

		// Create a new pflag set for each test, so flags don't bleed across.
		flagSet = pflag.NewFlagSet("test", pflag.ContinueOnError)

		// Save old CommandLine so we can restore it later.
		oldFlags = pflag.CommandLine
		pflag.CommandLine = flagSet

		// Re-init to register flags on our local FlagSet
		config.Init()
	})

	AfterEach(func() {
		// Restore original environment. We do this by clearing
		// everything and re-setting to origEnv
		os.Clearenv()
		for _, kv := range origEnv {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}

		// Reset log output.
		log.SetOutput(os.Stderr)

		// Restore execCommand.
		exec.Command = originalExecCommand

		// Restore pflag.CommandLine.
		pflag.CommandLine = oldFlags
	})

	Context("with no overrides at all", func() {
		It("uses the built-in defaults", func() {
			// We simulate that Git config is empty / fails.
			exec.Command = func(name string, arg ...string) exec.Cmd {
				return exec.NewMockCmd([]byte(""), fmt.Errorf("not found"))
			}

			// We do not parse flags or set environment variables.
			_ = flagSet.Parse([]string{})

			cfg, err := config.Load()
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Provider).To(Equal("openai"))
			Expect(cfg.Model).To(Equal("gpt-4o-mini"))
			Expect(cfg.OpenAIAPIKey).To(Equal(""))
		})
	})

	Context("when Git config provides values", func() {
		It("overrides defaults from Git for non-secret fields", func() {
			// Git returns a custom provider, e.g. "anthropic".
			exec.Command = func(name string, arg ...string) exec.Cmd {
				return exec.NewMockCmd([]byte("anthropic\n"), nil)
			}

			// We parse flags (but none given).
			_ = flagSet.Parse([]string{})

			cfg, err := config.Load()
			Expect(err).NotTo(HaveOccurred())

			// Because we re-call getGitConfigValue for model
			// it will be the same "anthropic" output. So let's
			// confirm both fields changed from default:
			Expect(cfg.Provider).To(Equal("anthropic"))
			Expect(cfg.Model).To(Equal("anthropic"))

			// Secret is not read from Git, remains default:
			Expect(cfg.OpenAIAPIKey).To(Equal(""))
		})

		It("logs errors if Git config fails but does not fail load", func() {
			exec.Command = func(name string, arg ...string) exec.Cmd {
				return exec.NewMockCmd([]byte(""), fmt.Errorf("git error"))
			}

			_ = flagSet.Parse([]string{})

			cfg, err := config.Load()
			Expect(err).NotTo(HaveOccurred())

			// We still get defaults even if Git config fails:
			Expect(cfg.Provider).To(Equal("openai"))
			Expect(cfg.Model).To(Equal("gpt-4o-mini"))

			// Check the log output for our error message:
			Expect(logBuf.String()).To(ContainSubstring("Error reading git config"))
		})
	})

	Context("when environment variables are set", func() {
		It("overrides defaults and Git config", func() {
			// Suppose Git says "anthropic".
			exec.Command = func(name string, arg ...string) exec.Cmd {
				return exec.NewMockCmd([]byte("anthropic\n"), nil)
			}

			// Also define environment variables for Provider, Model, and the secret:
			os.Setenv("GIT_AUTO_COMMIT_PROVIDER", "env-provider")
			os.Setenv("GIT_AUTO_COMMIT_MODEL", "env-model")
			os.Setenv("OPENAI_API_KEY", "env-secret")

			_ = flagSet.Parse([]string{})

			cfg, err := config.Load()
			Expect(err).NotTo(HaveOccurred())

			// Env beats Git:
			Expect(cfg.Provider).To(Equal("env-provider"))
			Expect(cfg.Model).To(Equal("env-model"))
			Expect(cfg.OpenAIAPIKey).To(Equal("env-secret"))
		})
	})

	Context("when flags are provided", func() {
		It("overrides everything else if non-empty", func() {
			// Git says "anthropic".
			exec.Command = func(name string, arg ...string) exec.Cmd {
				return exec.NewMockCmd([]byte("anthropic\n"), nil)
			}

			// Env says "env-provider".
			os.Setenv("GIT_AUTO_COMMIT_PROVIDER", "env-provider")
			os.Setenv("GIT_AUTO_COMMIT_MODEL", "env-model")
			os.Setenv("OPENAI_API_KEY", "env-secret")

			// Now we pass flags that override them all:
			err := flagSet.Parse([]string{
				"--provider=flag-provider",
				"--model=flag-model",
				"--openai-key=flag-secret",
			})
			Expect(err).NotTo(HaveOccurred())

			cfg, err := config.Load()
			Expect(err).NotTo(HaveOccurred())

			// pflags override environment and git.
			Expect(cfg.Provider).To(Equal("flag-provider"))
			Expect(cfg.Model).To(Equal("flag-model"))
			Expect(cfg.OpenAIAPIKey).To(Equal("flag-secret"))
		})

		It("does not override if the flag is empty", func() {
			// Suppose Git says "anthropic", environment says "env-provider"
			exec.Command = func(name string, arg ...string) exec.Cmd {
				return exec.NewMockCmd([]byte("anthropic\n"), nil)
			}
			os.Setenv("GIT_AUTO_COMMIT_PROVIDER", "env-provider")

			// Passing an empty flag to provider
			err := flagSet.Parse([]string{"--provider="})
			Expect(err).NotTo(HaveOccurred())

			cfg, err := config.Load()
			Expect(err).NotTo(HaveOccurred())

			// Because the flag was empty, environment overrides Git
			// and remains the final value:
			Expect(cfg.Provider).To(Equal("env-provider"))
		})
	})
})
