package exec

import "io"

// MockCmd is a stubbed command that you can use in tests. You can customize the
// behavior by setting fields or implementing side effects.
type MockCmd struct {
	// Simulated input for SetStdin().
	Stdin io.Reader
	// Simulated output for SetStdout().
	Stdout io.Writer
	// Simulated output for SetStderr().
	Stderr io.Writer

	// Simulated output for Output().
	MockOutput []byte

	// Simulated error for Output().
	MockError error

	// Track if methods were called (useful for verifying behavior in tests).
	OutputCalled    bool
	SetStdinCalled  bool
	SetStdoutCalled bool
	SetStderrCalled bool
	RunCalled       bool
}

// Ensure MockCommand satisfies the Command interface at compile time.
var _ Cmd = (*MockCmd)(nil)

// NewMockCmd is a convenience function that returns a *MockCommand. You can
// override exec.Command in your tests to return this.
func NewMockCmd(mockOutput []byte, mockError error) Cmd {
	return &MockCmd{
		MockOutput: mockOutput,
		MockError:  mockError,
	}
}

func (m *MockCmd) Output() ([]byte, error) {
	m.OutputCalled = true
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockOutput, nil
}

func (m *MockCmd) SetStdin(stdin io.Reader) {
	m.SetStdinCalled = true
	m.Stdin = stdin
}

func (m *MockCmd) SetStdout(stdout io.Writer) {
	m.SetStdoutCalled = true
	m.Stdout = stdout
}

func (m *MockCmd) SetStderr(stderr io.Writer) {
	m.SetStderrCalled = true
	m.Stderr = stderr
}

func (m *MockCmd) Run() error {
	m.RunCalled = true
	return m.MockError
}
