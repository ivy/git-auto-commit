package exec_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ivy/git-auto-commit/util/exec"
)

var _ = Describe("RealCmd", func() {
	Context("Output()", func() {
		It("returns the command's output when the command succeeds", func() {
			// Here we test the real command. We pick `echo hello` for simplicity.
			cmd := exec.NewRealCmd("echo", "hello")
			output, err := cmd.Output()
			Expect(err).NotTo(HaveOccurred())
			// `echo hello` adds a newline. Trim it if you prefer.
			Expect(strings.TrimSpace(string(output))).To(Equal("hello"))
		})

		It("returns an error if the command fails", func() {
			// Attempt to run a clearly invalid command (or pass invalid args)
			cmd := exec.NewRealCmd("this-command-does-not-exist")
			_, err := cmd.Output()
			Expect(err).To(HaveOccurred())
		})
	})
})
