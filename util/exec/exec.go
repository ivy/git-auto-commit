package exec

// Exec is the constructor for a Command. By default, it creates a RealCmd. You
// can override this with your own function (or inject a different constructor)
// in tests.
var Command = NewRealCmd

// Command is an interface encapsulating the methods of *exec.Cmd that are commonly used.
// You can add more methods here if you need more functionality (e.g., SetStdout, SetStderr, etc.).
type Cmd interface {
	Output() ([]byte, error)
}
