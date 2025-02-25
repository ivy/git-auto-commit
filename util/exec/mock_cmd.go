package exec

// MockCmd is a stubbed command that you can use in tests. You can customize the
// behavior by setting fields or implementing side effects.
type MockCmd struct {
	// Simulated output for Output().
	MockOutput []byte

	// Simulated error for Output().
	MockError error

	// Track if methods were called (useful for verifying behavior in tests).
	OutputCalled bool
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
