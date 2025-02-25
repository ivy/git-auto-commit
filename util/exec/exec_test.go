package exec_test

import (
	"testing"

	"github.com/ivy/git-auto-commit/util/exec"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Exec Suite")
}

var _ = Describe("Exec variable", func() {
	It("can be overridden to return a mock command", Serial, func() {
		originalExec := exec.Command
		defer func() {
			// Restore the original constructor after the test
			exec.Command = originalExec
		}()

		// Override Exec with a function that always returns a mock
		exec.Command = func(name string, args ...string) exec.Cmd {
			return exec.NewMockCmd([]byte("fake output"), nil)
		}

		cmd := exec.Command("whatever", "args")
		output, err := cmd.Output()
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(Equal([]byte("fake output")))

		mockCmd := cmd.(*exec.MockCmd)
		Expect(mockCmd.OutputCalled).To(BeTrue())
	})
})
