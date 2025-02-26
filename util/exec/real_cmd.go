package exec

import (
	"io"
	"os/exec"

	"github.com/ivy/git-auto-commit/util/log"
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
	log.Debugw("running command", "command", r.cmd.Path, "args", r.cmd.Args)
	return r.cmd.Output()
}

func (r *RealCmd) SetStdin(stdin io.Reader) {
	r.cmd.Stdin = stdin
}

func (r *RealCmd) SetStdout(stdout io.Writer) {
	r.cmd.Stdout = stdout
}

func (r *RealCmd) SetStderr(stderr io.Writer) {
	r.cmd.Stderr = stderr
}

func (r *RealCmd) Run() error {
	log.Debugw("running command", "command", r.cmd.Path, "args", r.cmd.Args)
	return r.cmd.Run()
}
