package git_auto_commit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/ivy/git-auto-commit/template"
	"github.com/ivy/git-auto-commit/util/exec"
	"github.com/ivy/git-auto-commit/util/git"
	"github.com/ivy/git-auto-commit/util/log"
)

// generatePRTitle generates a pull request title based on the supplied Git log.
func generatePRTitle(
	ctx context.Context, cfg *Config, description string,
) (string, error) {
	prompt, err := template.RenderString("prompt/pr_title.tmpl", map[string]any{
		"Description": description,
	})
	if err != nil {
		log.Errorw("failed to render pull request title template", "error", err)
		return "", err
	}

	client := openai.NewClient(
		option.WithAPIKey(cfg.OpenAIAPIKey),
	)

	stream := client.Chat.Completions.NewStreaming(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			}),
			Seed:  openai.Int(0),
			Model: openai.F(openai.ChatModel(cfg.Model)),
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

func generatePRDescription(ctx context.Context, cfg *Config) (string, error) {
	gitLog, err := git.Log()
	if err != nil {
		return "", fmt.Errorf("failed to get diff for PR: %w", err)
	}

	format, err := template.RenderString("format/pull_request.tmpl", nil)
	if err != nil {
		return "", err
	}

	prompt, err := template.RenderString(
		"prompt/pr_description.tmpl",
		map[string]any{
			"GitLog": gitLog,
			"Format": format,
		},
	)
	if err != nil {
		log.Errorw("failed to render pull request title template", "error", err)
		return "", err
	}

	client := openai.NewClient(
		option.WithAPIKey(cfg.OpenAIAPIKey),
	)

	stream := client.Chat.Completions.NewStreaming(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			}),
			Seed:  openai.Int(0),
			Model: openai.F(openai.ChatModel(cfg.Model)),
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

func AutoPullRequest(ctx context.Context, cfg *Config) error {
	log.Infow("starting auto-pr process",
		"verbose", cfg.Verbose,
		"extra_args", cfg.ExtraArgs)

	// 1. Generate a proposed PR description.
	prDescription, err := generatePRDescription(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to generate PR description: %w", err)
	}

	// 2. Generate a proposed PR title.
	prTitle, err := generatePRTitle(ctx, cfg, prDescription)
	if err != nil {
		return fmt.Errorf("failed to generate PR title: %w", err)
	}

	// Write the PR message to a temp file
	tempDir, err := os.MkdirTemp("", "git-auto-pr-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	f, err := os.Create(filepath.Join(tempDir, "PULLREQ_EDITMSG"))
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name())

	// 4. Invoke "gh pr create" with the final title and body.
	//    Optionally pass in base/head branches via cfg.ExtraArgs or explicitly.
	args := []string{
		"pr", "create",
		"--title", prTitle,
		"--body-file", f.Name(),
		"--web",
	}
	args = append(args, cfg.ExtraArgs...)

	log.Infow("creating pull request", "title", prTitle)

	cmd := exec.Command("gh", args...)
	cmd.SetStdin(os.Stdin)
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)
	return cmd.Run()
}
