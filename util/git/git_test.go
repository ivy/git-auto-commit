package git_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ivy/git-auto-commit/util/exec"
	"github.com/ivy/git-auto-commit/util/git"
)

func TestExec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Git Suite")
}

var _ = Describe("Status", func() {
	var (
		originalCommand func(name string, args ...string) exec.Cmd
	)

	BeforeEach(func() {
		// Save the original exec.Command function so we can restore it later
		originalCommand = exec.GetCommand()
	})

	AfterEach(func() {
		// Restore the original exec.Command after each test
		exec.SetCommand(originalCommand)
	})

	Context("when 'git status' succeeds", func() {
		It("returns the status output and no error", func() {
			// Create a mock command that simulates a successful output
			mockCmd := exec.NewMockCmd(
				[]byte("On branch main\nnothing to commit"), // Mocked output
				nil, // No error
			)

			// Override exec.Command to return our mock command
			exec.SetCommand(func(name string, args ...string) exec.Cmd {
				return mockCmd
			})

			status, err := git.Status()

			Expect(err).NotTo(HaveOccurred())
			// Verify that the mock's Output method was called
			Expect(mockCmd.(*exec.MockCmd).OutputCalled).To(BeTrue())
			Expect(status).To(Equal("On branch main\nnothing to commit"))
		})
	})

	Context("when 'git status' fails", func() {
		It("returns an error and empty output", func() {
			// Create a mock command that simulates an error
			mockCmd := exec.NewMockCmd(nil, fmt.Errorf("some error"))

			// Override exec.Command to return our mock command
			exec.SetCommand(func(name string, args ...string) exec.Cmd {
				return mockCmd
			})

			status, err := git.Status()

			Expect(err).To(HaveOccurred())
			// Verify that the mock's Output method was called
			Expect(mockCmd.(*exec.MockCmd).OutputCalled).To(BeTrue())
			Expect(status).To(BeEmpty())
		})
	})
})

var _ = Describe("Diff", func() {
	var (
		originalCommand func(name string, args ...string) exec.Cmd
	)

	BeforeEach(func() {
		// Save the original exec.Command function so we can restore it later
		originalCommand = exec.GetCommand()
	})

	AfterEach(func() {
		// Restore the original exec.Command after each test
		exec.SetCommand(originalCommand)
	})

	Context("when cached is true", func() {
		Context("when 'git diff --cached' succeeds", func() {
			It("returns the diff output and no error", func() {
				// Create a mock command that simulates a successful output
				mockCmd := exec.NewMockCmd(
					[]byte("diff --cached output"), // Mocked output
					nil,                            // No error
				)

				// Override exec.Command to return our mock command
				exec.SetCommand(func(name string, args ...string) exec.Cmd {
					return mockCmd
				})

				diff, err := git.Diff(true)

				Expect(err).NotTo(HaveOccurred())
				// Verify that the mock's Output method was called
				Expect(mockCmd.(*exec.MockCmd).OutputCalled).To(BeTrue())
				Expect(diff).To(Equal("diff --cached output"))
			})
		})

		Context("when 'git diff --cached' fails", func() {
			It("returns an error and empty output", func() {
				// Create a mock command that simulates an error
				mockCmd := exec.NewMockCmd(nil, fmt.Errorf("some error"))

				// Override exec.Command to return our mock command
				exec.SetCommand(func(name string, args ...string) exec.Cmd {
					return mockCmd
				})

				diff, err := git.Diff(true)

				Expect(err).To(HaveOccurred())
				// Verify that the mock's Output method was called
				Expect(mockCmd.(*exec.MockCmd).OutputCalled).To(BeTrue())
				Expect(diff).To(BeEmpty())
			})
		})
	})

	Context("when cached is false", func() {
		Context("when 'git diff' succeeds", func() {
			It("returns the diff output and no error", func() {
				// Create a mock command that simulates a successful output
				mockCmd := exec.NewMockCmd(
					[]byte("diff output"), // Mocked output
					nil,                   // No error
				)

				// Override exec.Command to return our mock command
				exec.SetCommand(func(name string, args ...string) exec.Cmd {
					return mockCmd
				})

				diff, err := git.Diff(false)

				Expect(err).NotTo(HaveOccurred())
				// Verify that the mock's Output method was called
				Expect(mockCmd.(*exec.MockCmd).OutputCalled).To(BeTrue())
				Expect(diff).To(Equal("diff output"))
			})
		})

		Context("when 'git diff' fails", func() {
			It("returns an error and empty output", func() {
				// Create a mock command that simulates an error
				mockCmd := exec.NewMockCmd(nil, fmt.Errorf("some error"))

				// Override exec.Command to return our mock command
				exec.SetCommand(func(name string, args ...string) exec.Cmd {
					return mockCmd
				})

				diff, err := git.Diff(false)

				Expect(err).To(HaveOccurred())
				// Verify that the mock's Output method was called
				Expect(mockCmd.(*exec.MockCmd).OutputCalled).To(BeTrue())
				Expect(diff).To(BeEmpty())
			})
		})
	})
})
