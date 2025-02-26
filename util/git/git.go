package git

import "github.com/ivy/git-auto-commit/util/exec"

// Status returns the output of `git status` command. It returns the status as a
// string and an error if the command fails.
func Status() (string, error) {
	cmd := exec.Command("git", "status")
	out, err := cmd.Output()
	return string(out), err
}

// Diff returns the output of `git diff` command. If cached is true, it returns
// the output of `git diff --cached`. It returns the diff as a string and an
// error if the command fails.
func Diff(cached bool) (string, error) {
	args := []string{"diff"}
	if cached {
		args = append(args, "--cached")
	}
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	return string(out), err
}
