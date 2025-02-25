package exec_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ivy/git-auto-commit/util/exec"
)

var _ = Describe("MockCmd", func() {
	Context("Output()", func() {
		It("returns the mocked output when there is no error", func() {
			mock := exec.NewMockCmd([]byte("mock output"), nil)
			output, err := mock.Output()
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal([]byte("mock output")))

			// Check that Output() was called, if you want to verify call counts:
			mockCmd := mock.(*exec.MockCmd)
			Expect(mockCmd.OutputCalled).To(BeTrue())
		})

		It("returns the mocked error and no output if an error is set", func() {
			mockErr := errors.New("mock error")
			mock := exec.NewMockCmd(nil, mockErr)
			output, err := mock.Output()
			Expect(err).To(MatchError("mock error"))
			Expect(output).To(BeNil())

			mockCmd := mock.(*exec.MockCmd)
			Expect(mockCmd.OutputCalled).To(BeTrue())
		})
	})
})
