package answerflow

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := struct {
		name  string
		input string
		want  *Client
	}{
		name:  "creates client with api key",
		input: "test-api-key",
		want: &Client{
			APIKey:  "test-api-key",
			BaseURL: "https://www.answeroverflow.com",
		},
	}

	t.Run(tests.name, func(t *testing.T) {
		got := New(tests.input)
		if got.APIKey != tests.want.APIKey {
			t.Errorf("APIKey = %q, want %q", got.APIKey, tests.want.APIKey)
		}
		if got.BaseURL != tests.want.BaseURL {
			t.Errorf("BaseURL = %q, want %q", got.BaseURL, tests.want.BaseURL)
		}
	})
}

func skipIfNoCreds(t *testing.T) {
	t.Helper()
	if os.Getenv("ANSWEROVERFLOW_API_KEY") == "" {
		t.Skip("skipping: ANSWEROVERFLOW_API_KEY not set")
	}
}

func TestCreateSolution(t *testing.T) {
	tests := []struct {
		name           string
		messageID      string
		solutionID     string
		server         *httptest.Server
		wantResponse   *CreateSolutionResponse
		wantErr        bool
		errContains    string
		isIntegration  bool
	}{
		{
			name:       "successful solution creation",
			messageID:  "msg123",
			solutionID: "sol456",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("method = %q, want %q", r.Method, http.MethodPost)
				}
				if !strings.Contains(r.URL.Path, "msg123") {
					t.Errorf("URL path = %q, want to contain %q", r.URL.Path, "msg123")
				}
				if r.Header.Get("x-api-key") != "test-key" {
					t.Errorf("x-api-key = %q, want %q", r.Header.Get("x-api-key"), "test-key")
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Content-Type = %q, want %q", r.Header.Get("Content-Type"), "application/json")
				}
				if r.Header.Get("User-Agent") == "" {
					t.Error("User-Agent header not set")
				}

				var req CreateSolutionRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("failed to decode request: %v", err)
				}
				if req.SolutionID != "sol456" {
					t.Errorf("SolutionID = %q, want %q", req.SolutionID, "sol456")
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(CreateSolutionResponse{Success: true})
			})),
			wantResponse: &CreateSolutionResponse{Success: true},
		},
		{
			name:       "server returns error",
			messageID:  "msg123",
			solutionID: "sol456",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(CreateSolutionResponse{
					Success: false,
					Error:   "invalid solution id",
				})
			})),
			wantResponse: &CreateSolutionResponse{
				Success: false,
				Error:   "invalid solution id",
			},
		},
		{
			name:       "invalid json response",
			messageID:  "msg123",
			solutionID: "sol456",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("not json"))
			})),
			wantErr:     true,
			errContains: "failed to decode response",
		},
		{
			name:       "server returns 500",
			messageID:  "msg123",
			solutionID: "sol456",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})),
			wantErr:     true,
			errContains: "failed to decode response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isIntegration {
				skipIfNoCreds(t)
			}
			if tt.server != nil {
				defer tt.server.Close()
			}

			client := &Client{
				APIKey:  "test-key",
				BaseURL: tt.server.URL,
			}

			got, err := client.CreateSolution(tt.messageID, tt.solutionID)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errContains)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %q, want to contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Success != tt.wantResponse.Success {
				t.Errorf("Success = %v, want %v", got.Success, tt.wantResponse.Success)
			}
			if got.Error != tt.wantResponse.Error {
				t.Errorf("Error = %q, want %q", got.Error, tt.wantResponse.Error)
			}
		})
	}
}
