package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/facebookgo/flagenv"
	claudecode "github.com/humanlayer/humanlayer/claudecode-go"

	_ "embed"

	_ "github.com/joho/godotenv/autoload"
)

var (
	anthropicAuthToken = flag.String("anthropic-auth-token", "hunter2", "Anthropic API token")
	anthropicBaseURL   = flag.String("anthropic-base-url", "http://localhost:11434", "Anthropic API base URL")
	anthropicModel     = flag.String("anthropic-model", "glm-4.7-flash:latest", "Anthropic AI model to use for all levels of agentic function")
	zhipuAPIKey        = flag.String("zhipu-api-key", "", "API key for z.ai (Zhipu)")

	//go:embed prompts/optimized.txt
	testPrompt string
)

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
	client, err := claudecode.NewClient()
	if err != nil {
		return fmt.Errorf("can't open Claude Code: %w", err)
	}

	cwd, _ := os.Getwd()

	sess, err := client.Launch(claudecode.SessionConfig{
		Query:        testPrompt,
		OutputFormat: claudecode.OutputStreamJSON,
		AllowedTools: []string{"mcp__*", "Bash(*)", "WebSearch", "Read", "Write", "Grep", "Glob"},
		// PermissionPromptTool: "mcp__approval__prompt-user",
		AdditionalDirectories: []string{filepath.Join(cwd, "var", "*")},
		Verbose:               true,
		WorkingDir:            cwd,

		MCPConfig: &claudecode.MCPConfig{
			MCPServers: map[string]claudecode.MCPServer{
				"approval": {
					Command: "~/go/bin/mcp-yolo-approval",
					//Command: "go",
					//Args:    []string{"run", "../mcp-yolo-approval"},
				},
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
					lg.Info("using tool", "tool", part.Name, "input", part.Input)
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

	return nil
}
