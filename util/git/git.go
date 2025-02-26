package git

import (
	"fmt"
	"strings"

	"github.com/ivy/git-auto-commit/util/exec"
)

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

// DefaultBranch returns the name of the default branch. It returns the branch
// name as a string and an error if the command fails.
func DefaultBranch() (string, error) {
	cmd := exec.Command("git", "remote", "show", "origin")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "HEAD branch:") {
			parts := strings.Split(line, ":")
			return strings.TrimSpace(parts[len(parts)-1]), nil
		}
	}

	return "", fmt.Errorf("default branch not found")
}

// Log returns the output of `git log --patch` command. It returns the log as a
// string and an error if the command fails.
func Log() (string, error) {
	defaultBranch, err := DefaultBranch()
	if err != nil {
		return "", err
	}

	cmd := exec.Command(
		"git", "log",
		fmt.Sprintf("%s/%s...%s", "origin", defaultBranch, "HEAD"),
	)
	out, err := cmd.Output()
	return string(out), err
}
