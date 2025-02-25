package exec

import "sync"

var (
	mu sync.RWMutex

	// Exec is the constructor for a Command. By default, it creates a RealCmd. You
	// can override this with your own function (or inject a different constructor)
	// in tests.
	command = NewRealCmd
)

// Command is an interface encapsulating the methods of *exec.Cmd that are commonly used.
// You can add more methods here if you need more functionality (e.g., SetStdout, SetStderr, etc.).
type Cmd interface {
	Output() ([]byte, error)
}

func Command(name string, arg ...string) Cmd {
	mu.RLock()
	defer mu.RUnlock()
	return command(name, arg...)
}

func GetCommand() func(string, ...string) Cmd {
	mu.RLock()
	defer mu.RUnlock()
	return command
}

func SetCommand(c func(string, ...string) Cmd) {
	mu.Lock()
	defer mu.Unlock()
	command = c
}
