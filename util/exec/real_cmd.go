package exec

import (
	"os/exec"
)

// RealCmd wraps an *exec.Cmd and implements the Cmd interface.
type RealCmd struct {
	cmd *exec.Cmd
}

// Ensure RealCommand satisfies the Command interface at compile time.
var _ Cmd = (*RealCmd)(nil)

func NewRealCmd(name string, arg ...string) Cmd {
	return &RealCmd{
		cmd: exec.Command(name, arg...),
	}
}

// Output wraps (*exec.Cmd).Output()
func (r *RealCmd) Output() ([]byte, error) {
	return r.cmd.Output()
}
