package main_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "git-auto-pr Suite")
}

var _ = Describe("Integration with main.go binary", func() {
	It("prints version and exits 0 when run with --version", func() {
		cmd := exec.Command("go", "run", "./main.go", "--version")

		stdoutBuf := &bytes.Buffer{}
		stderrBuf := &bytes.Buffer{}
		cmd.Stdout = stdoutBuf
		cmd.Stderr = stderrBuf

		err := cmd.Run()
		exitCode := 0
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}

		Expect(exitCode).To(Equal(0), "Expected exit code 0.")
		Expect(strings.TrimSpace(stdoutBuf.String())).To(ContainSubstring("git auto-pr"))
		Expect(stderrBuf.String()).To(BeEmpty())
	})
})
