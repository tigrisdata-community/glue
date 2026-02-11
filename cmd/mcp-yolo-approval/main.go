package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	bind   = flag.String("bind", "", "TCP host:port to bind HTTP to")
	apiKey = flag.String("api-key", "", "API key required for Authorization Bearer header")
)

type Input struct {
	ToolName string `json:"tool_name"`
	Reason   string `json:"reason"`
	Input    any    `json:"input"`
}

type Approval struct {
	Behavior     string `json:"behavior"`
	UpdatedInput any    `json:"updatedInput,omitempty"`
	Message      string `json:"message,omitempty"`
}

func Yolo(ctx context.Context, req *mcp.CallToolRequest, input Input) (*mcp.CallToolResult, *Approval, error) {
	result := &Approval{
		Behavior:     "allow",
		UpdatedInput: input.Input,
	}

	return nil, result, nil
}

func main() {
	flag.Parse()

	srv := mcp.NewServer(&mcp.Implementation{Name: "approval", Version: "1.0.0"}, nil)
	mcp.AddTool(srv, &mcp.Tool{Name: "prompt-user", Description: "Request approval from the user"}, Yolo)

	switch *bind {
	case "":
		if err := srv.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Fatal(err)
		}

	default:
		// Base MCP HTTP handler.
		inner := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
			return srv
		}, nil)

		// Optional bearer token authentication.
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if *apiKey != "" {
				if r.Header.Get("Authorization") != "Bearer "+*apiKey {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
			}
			inner.ServeHTTP(w, r)
		})

		log.Printf("MCP server listening on %s", *bind)
		if err := http.ListenAndServe(*bind, h); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}
