package git

import "github.com/ivy/git-auto-commit/util/exec"

// Status returns the output of `git status` command. It returns the status as a
// string and an error if the command fails.
func Status() (string, error) {
	cmd := exec.Command("git", "status")
	out, err := cmd.Output()
	return string(out), err
}
