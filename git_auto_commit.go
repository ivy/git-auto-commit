package git_auto_commit

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	stdexec "os/exec"
	"path/filepath"
	"strings"
	txtTemplate "text/template"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/ivy/git-auto-commit/config"
	"github.com/ivy/git-auto-commit/template"
	"github.com/ivy/git-auto-commit/util/exec"
	"github.com/ivy/git-auto-commit/util/git"
	"github.com/ivy/git-auto-commit/util/log"
)

const (
	scissors    = "------------------------ >8 ------------------------"
	commentChar = "#"
)

// Command is the configuration for the git-auto-commit and
// git-auto-pr commands.
type Config struct {
	*config.Config

	// Verbose opens an editor for the user to review messages.
	Verbose bool

	// Yes skips the editor and directly commits the message.
	Yes bool

	// Message provides additional context for the commit message. It's supplied
	// by the user on the command line.
	Message string

	// ExtraArgs are additional arguments to pass to the used git/gh command.
	ExtraArgs []string
}

var editorFallbacks = []string{"nano", "vim", "vi"}

// commitMessagePrompt is a template for the commit message prompt.
var commitMessagePrompt = txtTemplate.Must(txtTemplate.New("commitMessage").Parse(
	`You are a helpful assistant who generates commit messages for Git.

Commit messages follow this format:

{{.Format}}

---

The following changes have been staged for commit:

{{.Staged}}

---

Additional context for the commit message: {{.Message}}

---

Generate a commit message for the changes, following the format above.`,
))

// commitFooter is added to the end of verbose commits before opening them in
// the editor.
var commitFooter = txtTemplate.Must(txtTemplate.New("commitFooter").Parse(
	`Please edit the commit message to your liking. Lines starting
with '{{.CommentChar}}' will be ignored, and an empty message aborts the commit.

{{.GitStatus}}
------------------------ >8 ------------------------
Do not modify or remove the line above.
Everything below it will be ignored.`))

// prefixLines prefixes each line of the input reader with the given prefix
// string and writes the result to the output writer.  It returns an error if
// one occurs during reading or writing.
func prefixLines(r io.Reader, w io.Writer, prefix string) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		prefixedLine := prefix + line + "\n" // Add newline back
		_, err := w.Write([]byte(prefixedLine))
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

// GenerateCommitMessage generates a commit message for the given staged changes
// and Config using AI.
func GenerateCommitMessage(ctx context.Context, config *Config, staged string) (string, error) {
	log.Debugw("generating commit message",
		"model", config.Model,
		"message_context", config.Message)

	client := openai.NewClient(
		option.WithAPIKey(config.OpenAIAPIKey),
	)

	format, err := template.RenderString("format/commit.tmpl", nil)
	if err != nil {
		log.Errorw("failed to render commit message format",
			"error", err)
		return "", err
	}

	prompt := &bytes.Buffer{}
	err = commitMessagePrompt.Execute(prompt, map[string]any{
		"Staged":  staged,
		"Format":  format,
		"Message": config.Message,
	})
	if err != nil {
		log.Errorw("failed to execute commit message template",
			"error", err)
		return "", err
	}
	log.Debugw("commit message template executed", "prompt", prompt.String())

	stream := client.Chat.Completions.NewStreaming(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt.String()),
			}),
			Seed:  openai.Int(0),
			Model: openai.F(openai.ChatModel(config.Model)),
		},
	)

	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)
		log.Debugw("stream chunk received", "chunk", chunk)

		// if content, ok := acc.JustFinishedContent(); ok {
		// 	return content, nil
		// }

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			log.Warnw("AI refused to generate commit message",
				"refusal", refusal)
			return "", fmt.Errorf("refusal: %s", refusal)
		}
	}

	if err := stream.Err(); err != nil {
		log.Errorw("stream error while generating commit message",
			"error", err)
		return "", err
	}

	return acc.Choices[0].Message.Content, nil
}

// AutoCommit uses Git to commit staged changes, generating a commit message
// using AI.
func AutoCommit(ctx context.Context, config *Config) error {
	log.Infow("starting auto-commit process",
		"verbose", config.Verbose,
		"extra_args", config.ExtraArgs)

	// 1. Get the staged changes.
	// TODO(ivy): handle amending commits
	staged, err := git.Diff(true)
	if err != nil {
		log.Errorw("failed to get staged changes",
			"error", err)
		return err
	}

	// 2. Generate a commit message.
	message, err := GenerateCommitMessage(ctx, config, string(staged))
	if err != nil {
		log.Errorw("failed to generate commit message",
			"error", err)
		return err
	}

	log.Debugw("generated commit message",
		"message", message)

	// 3. Optionally, open the editor for the user to review the message.
	if config.Verbose {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			// If no editor is set, use the first installed fallback.
			for _, fallback := range editorFallbacks {
				if _, err := stdexec.LookPath(fallback); err == nil {
					editor = fallback
					break
				}
			}
		}
		if editor == "" {
			log.Warnw("no editor found",
				"fallbacks_tried", editorFallbacks)
			return errors.New("no editor found, set $EDITOR")
		}

		log.Infow("opening editor for commit message review",
			"editor", editor)

		tempDir, err := os.MkdirTemp("", "git-auto-commit-*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)

		// Write message to a temporary file
		f, err := os.Create(filepath.Join(tempDir, "COMMIT_EDITMSG"))
		if err != nil {
			return err
		}
		defer f.Close()
		defer os.Remove(f.Name())

		gitStatus, err := git.Status()
		if err != nil {
			return err
		}

		footer := new(bytes.Buffer)
		err = commitFooter.Execute(footer, map[string]string{
			// TODO(ivy): use core.commentChar
			"CommentChar": commentChar,
			"Scissors":    scissors,
			"GitStatus":   gitStatus,
		})
		if err != nil {
			return err
		}

		if _, err = f.WriteString(message + "\n\n"); err != nil {
			return err
		}

		r := bytes.NewBufferString(footer.String())
		w := new(bytes.Buffer)
		if err = prefixLines(r, w, commentChar+" "); err != nil {
			return err
		}
		if _, err = f.WriteString(w.String()); err != nil {
			return err
		}

		if _, err = f.WriteString(staged); err != nil {
			return err
		}

		// Open an editor to allow edits to the generated message.
		editorParts := append(strings.Split(editor, " "), f.Name())
		cmd := exec.Command(editorParts[0], editorParts[1:]...)
		cmd.SetStdin(os.Stdin)
		cmd.SetStdout(os.Stdout)
		cmd.SetStderr(os.Stderr)
		if err := cmd.Run(); err != nil {
			return err
		}

		args := []string{"commit", "--file", f.Name()}
		args = append(args, config.ExtraArgs...)

		// Open an editor to confirm the commit with the generated message.
		cmd = exec.Command("git", args...)
		cmd.SetStdin(os.Stdin)
		cmd.SetStdout(os.Stdout)
		cmd.SetStderr(os.Stderr)
		return cmd.Run()
	}

	log.Infow("committing changes",
		"extra_args", config.ExtraArgs)

	// 3. Otherwise, commit the changes and pass any extra args.
	cmd := exec.Command(
		"git",
		append([]string{"commit", "--file", "-"}, config.ExtraArgs...)...,
	)
	cmd.SetStdin(strings.NewReader(message))
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)
	return cmd.Run()
}
