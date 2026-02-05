package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/facebookgo/flagenv"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	claudecode "github.com/humanlayer/humanlayer/claudecode-go"

	_ "github.com/joho/godotenv/autoload"
)

var (
	anthropicAuthToken = flag.String("anthropic-auth-token", "hunter2", "Anthropic API token")
	anthropicBaseURL   = flag.String("anthropic-base-url", "http://localhost:11434", "Anthropic API base URL")
	anthropicModel     = flag.String("anthropic-model", "glm-4.7-flash:latest", "Anthropic AI model to use for all levels of agentic function")
	zhipuAPIKey        = flag.String("zhipu-api-key", "", "API key for z.ai (Zhipu)")
	outputFolder       = flag.String("output-folder", "./var", "Output folder for generated content")

	//go:embed prompts/*.tmpl.txt
	prompts embed.FS

	ErrNoInputTopic   = errors.New("no topic defined")
	ErrNoRelevantDocs = errors.New("no relevant documentation defined")
	NoChangesToCommit = errors.New("no changes to commit")
)

type Input struct {
	Topic        string `json:"topic"`
	RelevantDocs string `json:"relevantDocs"`
	OutputFolder string `json:"outputFolder"`
}

func (i Input) Valid() error {
	var errs []error

	if i.Topic == "" {
		errs = append(errs, ErrNoInputTopic)
	}

	if i.RelevantDocs == "" {
		errs = append(errs, ErrNoRelevantDocs)
	}

	if len(errs) != 0 {
		return fmt.Errorf("Input failed validation: %w", errors.Join(errs...))
	}

	return nil
}

func commitChanges(repoPath, outputFolder, topic string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("can't open git repo: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("can't get worktree: %w", err)
	}

	// Get status to find changes in output folder
	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("can't get git status: %w", err)
	}

	// Collect files in output folder that have changes
	var filesToCommit []string
	outputFolderAbs, err := filepath.Abs(outputFolder)
	if err != nil {
		return fmt.Errorf("can't get absolute path for output folder: %w", err)
	}

	for file, st := range status {
		if st.Worktree == git.Unmodified {
			continue
		}
		// Check if file is in or under the output folder
		absPath, err := filepath.Abs(file)
		if err != nil {
			continue
		}
		if strings.HasPrefix(absPath, outputFolderAbs) || strings.HasPrefix(file, outputFolder) {
			filesToCommit = append(filesToCommit, file)
		}
	}

	if len(filesToCommit) == 0 {
		return NoChangesToCommit
	}

	// Add all changed files in output folder
	for _, file := range filesToCommit {
		if _, err := worktree.Add(file); err != nil {
			return fmt.Errorf("can't stage file %s: %w", file, err)
		}
	}

	// Create commit message
	commitMsg := fmt.Sprintf("docs: %s\n\nAssisted-by: GLM 4.7 via Claude Code\nSigned-off-by: Xe Iaso <xe@tigrisdata.com>", topic)

	// Commit
	_, err = worktree.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Xe Iaso",
			Email: "xe@tigrisdata.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("can't commit: %w", err)
	}

	slog.Info("committed changes", "output_folder", outputFolder, "files", len(filesToCommit), "topic", topic)
	return nil
}

func main() {
	flagenv.Parse()
	flag.Parse()

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errors.New("main exited"))

	slog.Info(
		"starting up",
		"has-anthropic-api-token", *anthropicAuthToken != "",
		"anthropic-base-url", *anthropicBaseURL,
		"anthropic-model", *anthropicModel,
	)

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
	}
}

func run(ctx context.Context) error {
	var input Input

	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		return fmt.Errorf("can't read input JSON: %w", err)
	}

	if err := input.Valid(); err != nil {
		return fmt.Errorf("can't validate input JSON: %w", err)
	}

	input.OutputFolder = *outputFolder

	var promptBuilder strings.Builder

	tmpl, err := template.ParseFS(prompts, "prompts/*.tmpl.txt")
	if err != nil {
		return fmt.Errorf("can't parse templates: %w", err)
	}
	if err := tmpl.ExecuteTemplate(&promptBuilder, "optimized.tmpl.txt", input); err != nil {
		return fmt.Errorf("can't hydrate prompt: %w", err)
	}

	client, err := claudecode.NewClient()
	if err != nil {
		return fmt.Errorf("can't open Claude Code: %w", err)
	}

	cwd, _ := os.Getwd()

	sess, err := client.Launch(claudecode.SessionConfig{
		Query:        promptBuilder.String(),
		OutputFormat: claudecode.OutputStreamJSON,
		AllowedTools: []string{"mcp__web-reader__*", "mcp__tigris-discord__*", "Bash(*)", "Bash(find*)", "WebSearch", "Read", "Write", "Grep", "Glob", "Edit", "Update"},
		// PermissionPromptTool: "mcp__approval__prompt-user",
		AdditionalDirectories: []string{*outputFolder, cwd},
		Verbose:               true,
		WorkingDir:            cwd,

		MCPConfig: &claudecode.MCPConfig{
			MCPServers: map[string]claudecode.MCPServer{
				"tigris-discord": {
					Type: "http",
					URL:  "https://community.tigrisdata.com/mcp",
				},
				"web-reader": {
					Type: "http",
					URL:  "https://api.z.ai/api/mcp/web_reader/mcp",
					Headers: map[string]string{
						"Authorization": "Bearer " + *zhipuAPIKey,
					},
				},
			},
		},

		Env: map[string]string{
			"ANTHROPIC_AUTH_TOKEN":           *anthropicAuthToken,
			"ANTHROPIC_BASE_URL":             *anthropicBaseURL,
			"ANTHROPIC_DEFAULT_HAIKU_MODEL":  *anthropicModel,
			"ANTHROPIC_DEFAULT_SONNET_MODEL": *anthropicModel,
			"ANTHROPIC_DEFAULT_OPUS_MODEL":   *anthropicModel,
		},
	})
	if err != nil {
		return fmt.Errorf("can't open Claude Code session: %w", err)
	}

	lg := slog.With()
	lg.Info("started agent")

	for event := range sess.Events {
		lg.Info("got event", "session", sess.ID, "type", event.Type, "subtype", event.Subtype, "is_error", event.IsError)
		if event.IsError {
			lg.Error("execution error", "err", event.Error)
		}
		if event.Message != nil {
			for _, part := range event.Message.Content {
				switch part.Type {
				case "tool_use":
					switch part.Name {
					case "TodoWrite":
						inputMap := make(map[string]any, len(part.Input))
						for k, v := range part.Input {
							inputMap[k] = v
						}
						if todoList, err := ParseTodoFromMap(inputMap); err == nil {
							for _, todo := range todoList.Todos {
								lg.Info("todo", "status", todo.Status, "content", todo.Content)
							}
						} else {
							lg.Info("using tool", "tool", part.Name, "input", part.Input)
						}
					case "mcp__web-reader__webReader":
						if url, ok := part.Input["url"].(string); ok {
							lg.Info("fetching docs", "url", url)
						} else {
							lg.Info("using tool", "tool", part.Name, "input", part.Input)
						}
					case "Read":
						if path, ok := part.Input["file_path"].(string); ok {
							lg.Info("reading file", "path", path)
						} else {
							lg.Info("using tool", "tool", part.Name, "input", part.Input)
						}
					case "Write":
						if path, ok := part.Input["file_path"].(string); ok {
							lg.Info("writing file", "path", path)
						} else {
							lg.Info("using tool", "tool", part.Name, "input", part.Input)
						}
					case "Edit":
						if path, ok := part.Input["file_path"].(string); ok {
							lg.Info("editing file", "path", path)
						} else {
							lg.Info("using tool", "tool", part.Name, "input", part.Input)
						}
					default:
						lg.Info("using tool", "tool", part.Name, "input", part.Input)
					}
				case "tool_result":
					lg.Info("tool result", "tool", part.Name)
					json.NewEncoder(os.Stdout).Encode(part)
					fmt.Println()
				}
				if part.Text != "" {
					fmt.Printf("%s> %s\n", event.Message.Role, part.Text)
				}
				lg.Info("message type", "session", sess.ID, "type", event.Type, "subtype", event.Subtype, "message_type", part.Type)
			}
		}
	}

	result, err := sess.Wait()
	if err != nil {
		return fmt.Errorf("can't get session result: %w", err)
	}

	if result.IsError {
		lg.Error("got error result", "session", sess.ID, "type", result.Type, "subtype", result.Subtype, "cost_usd", result.CostUSD, "duration_ms", result.DurationMS, "num_turns", result.NumTurns, "err", result.Error)
		return fmt.Errorf("got error from Claude: %s", result.Error)
	} else {
		lg.Info("got result", "session", sess.ID, "type", result.Type, "subtype", result.Subtype, "cost_usd", result.CostUSD, "duration_ms", result.DurationMS, "num_turns", result.NumTurns)
	}

	// Commit any changes made in the output folder
	if err := commitChanges(*outputFolder, *outputFolder, input.Topic); err != nil {
		if errors.Is(err, NoChangesToCommit) {
			slog.Info("no changes to commit in output folder")
		} else {
			slog.Error("failed to commit changes", "err", err)
		}
	}

	return nil
}
