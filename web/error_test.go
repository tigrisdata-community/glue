package web

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		name         string
		wantStatus   int
		responseCode int
		responseBody string
		wantErr      bool
		checkError   func(*testing.T, error)
	}{
		{
			name:         "creates error for 404 response",
			wantStatus:   http.StatusOK,
			responseCode: http.StatusNotFound,
			responseBody: "not found",
			wantErr:      true,
			checkError: func(t *testing.T, err error) {
				webErr, ok := err.(*Error)
				if !ok {
					t.Fatalf("expected *Error, got %T", err)
				}
				if webErr.WantStatus != http.StatusOK {
					t.Errorf("WantStatus = %d, want %d", webErr.WantStatus, http.StatusOK)
				}
				if webErr.GotStatus != http.StatusNotFound {
					t.Errorf("GotStatus = %d, want %d", webErr.GotStatus, http.StatusNotFound)
				}
				if webErr.ResponseBody != "not found" {
					t.Errorf("ResponseBody = %s, want 'not found'", webErr.ResponseBody)
				}
			},
		},
		{
			name:         "creates error for 500 response",
			wantStatus:   http.StatusOK,
			responseCode: http.StatusInternalServerError,
			responseBody: "internal server error",
			wantErr:      true,
			checkError: func(t *testing.T, err error) {
				webErr, ok := err.(*Error)
				if !ok {
					t.Fatalf("expected *Error, got %T", err)
				}
				if webErr.GotStatus != http.StatusInternalServerError {
					t.Errorf("GotStatus = %d, want %d", webErr.GotStatus, http.StatusInternalServerError)
				}
			},
		},
		{
			name:         "creates error with POST method",
			wantStatus:   http.StatusCreated,
			responseCode: http.StatusNotFound,
			responseBody: "bad request",
			wantErr:      true,
			checkError: func(t *testing.T, err error) {
				webErr, ok := err.(*Error)
				if !ok {
					t.Fatalf("expected *Error, got %T", err)
				}
				if webErr.Method != http.MethodPost {
					t.Errorf("Method = %s, want POST", webErr.Method)
				}
			},
		},
		{
			name:         "creates error with URL",
			wantStatus:   http.StatusOK,
			responseCode: http.StatusNotFound,
			responseBody: "error",
			wantErr:      true,
			checkError: func(t *testing.T, err error) {
				webErr, ok := err.(*Error)
				if !ok {
					t.Fatalf("expected *Error, got %T", err)
				}
				if webErr.URL == nil {
					t.Fatal("URL is nil")
				}
				if !strings.HasPrefix(webErr.URL.Scheme, "http") {
					t.Errorf("URL scheme = %s, want http(s)", webErr.URL.Scheme)
				}
			},
		},
		{
			name:         "creates error with empty response body",
			wantStatus:   http.StatusOK,
			responseCode: http.StatusNotFound,
			responseBody: "",
			wantErr:      true,
			checkError: func(t *testing.T, err error) {
				webErr, ok := err.(*Error)
				if !ok {
					t.Fatalf("expected *Error, got %T", err)
				}
				if webErr.ResponseBody != "" {
					t.Errorf("ResponseBody = %s, want empty", webErr.ResponseBody)
				}
			},
		},
		{
			name:         "creates error with large response body",
			wantStatus:   http.StatusOK,
			responseCode: http.StatusNotFound,
			responseBody: strings.Repeat("x", 10000),
			wantErr:      true,
			checkError: func(t *testing.T, err error) {
				webErr, ok := err.(*Error)
				if !ok {
					t.Fatalf("expected *Error, got %T", err)
				}
				if len(webErr.ResponseBody) != 10000 {
					t.Errorf("ResponseBody length = %d, want 10000", len(webErr.ResponseBody))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			var resp *http.Response
			var err error

			if tt.name == "creates error with POST method" {
				resp, err = http.Post(server.URL, "application/json", strings.NewReader("{}"))
			} else {
				resp, err = http.Get(server.URL)
			}

			if err != nil {
				t.Fatalf("failed to create test response: %v", err)
			}

			err = NewError(tt.wantStatus, resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewError() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.checkError != nil && err != nil {
				tt.checkError(t, err)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      Error
		contains []string
	}{
		{
			name: "error string contains method",
			err: Error{
				Method: "GET",
				URL:    &url.URL{Scheme: "http", Host: "localhost"},
			},
			contains: []string{"GET"},
		},
		{
			name: "error string contains URL",
			err: Error{
				Method: "GET",
				URL:    &url.URL{Scheme: "https", Host: "example.com", Path: "/test"},
			},
			contains: []string{"GET", "example.com"},
		},
		{
			name: "error string contains status codes",
			err: Error{
				Method:     "POST",
				URL:        &url.URL{Scheme: "http", Host: "localhost"},
				WantStatus: 200,
				GotStatus:  404,
			},
			contains: []string{"POST", "200", "404"},
		},
		{
			name: "error string contains response body",
			err: Error{
				Method:       "GET",
				URL:          &url.URL{Scheme: "http", Host: "localhost"},
				WantStatus:   200,
				GotStatus:    500,
				ResponseBody: "internal error",
			},
			contains: []string{"200", "500", "internal error"},
		},
		{
			name: "error string with all fields",
			err: Error{
				Method:       "DELETE",
				URL:          &url.URL{Scheme: "https", Host: "api.example.com", Path: "/resource/123"},
				WantStatus:   http.StatusNoContent,
				GotStatus:    http.StatusForbidden,
				ResponseBody: "access denied",
			},
			contains: []string{"DELETE", "api.example.com", "204", "403", "access denied"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.err.Error()
			for _, substr := range tt.contains {
				if !strings.Contains(errStr, substr) {
					t.Errorf("Error() string = %s, does not contain %s", errStr, substr)
				}
			}
		})
	}
}

func TestError_LogValue(t *testing.T) {
	// Helper URL for tests that don't care about the specific URL
	testURL := &url.URL{Scheme: "http", Host: "localhost"}

	tests := []struct {
		name       string
		err        Error
		wantKeys   []string
		wantValues map[string]any
	}{
		{
			name: "LogValue contains want_status",
			err: Error{
				URL:        testURL,
				WantStatus: 200,
			},
			wantKeys: []string{"want_status"},
			wantValues: map[string]any{
				"want_status": int64(200),
			},
		},
		{
			name: "LogValue contains got_status",
			err: Error{
				URL:       testURL,
				GotStatus: 404,
			},
			wantKeys: []string{"got_status"},
			wantValues: map[string]any{
				"got_status": int64(404),
			},
		},
		{
			name: "LogValue contains url",
			err: Error{
				URL: &url.URL{Scheme: "https", Host: "example.com", Path: "/test"},
			},
			wantKeys: []string{"url"},
			wantValues: map[string]any{
				"url": "https://example.com/test",
			},
		},
		{
			name: "LogValue contains method",
			err: Error{
				URL:    testURL,
				Method: "POST",
			},
			wantKeys: []string{"method"},
			wantValues: map[string]any{
				"method": "POST",
			},
		},
		{
			name: "LogValue contains body",
			err: Error{
				URL:          testURL,
				ResponseBody: "error body",
			},
			wantKeys: []string{"body"},
			wantValues: map[string]any{
				"body": "error body",
			},
		},
		{
			name: "LogValue with all fields",
			err: Error{
				WantStatus:   200,
				GotStatus:    500,
				URL:          &url.URL{Scheme: "https", Host: "api.test.com", Path: "/endpoint"},
				Method:       "PUT",
				ResponseBody: "server error",
			},
			wantKeys: []string{"want_status", "got_status", "url", "method", "body"},
			wantValues: map[string]any{
				"want_status": int64(200),
				"got_status":  int64(500),
				"url":         "https://api.test.com/endpoint",
				"method":      "PUT",
				"body":        "server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := tt.err.LogValue()

			// Check that val is a GroupValue
			if val.Kind() != slog.KindGroup {
				t.Errorf("LogValue() kind = %v, want KindGroup", val.Kind())
				return
			}

			// Get the group attributes
			attrs := val.Group()

			// Create a map of attribute values
			attrMap := make(map[string]any)
			for _, attr := range attrs {
				attrMap[attr.Key] = attr.Value.Any()
			}

			// Check that all expected keys are present
			for _, key := range tt.wantKeys {
				if _, ok := attrMap[key]; !ok {
					t.Errorf("LogValue() missing key %s. Got keys: %v", key, attrMap)
				}
			}

			// Check specific values
			for key, wantVal := range tt.wantValues {
				gotVal, ok := attrMap[key]
				if !ok {
					t.Errorf("LogValue() missing key %s", key)
					continue
				}
				if gotVal != wantVal {
					t.Errorf("LogValue() key %s = %v, want %v", key, gotVal, wantVal)
				}
			}
		})
	}
}

func TestError_Integration(t *testing.T) {
	// Integration test that simulates a real error flow
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("expected /test path, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("resource not found"))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL + "/test")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}

	err = NewError(http.StatusOK, resp)
	if err == nil {
		t.Fatal("NewError() returned nil, expected error")
	}

	webErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error type, got %T", err)
	}

	// Verify all fields
	if webErr.WantStatus != http.StatusOK {
		t.Errorf("WantStatus = %d, want %d", webErr.WantStatus, http.StatusOK)
	}
	if webErr.GotStatus != http.StatusNotFound {
		t.Errorf("GotStatus = %d, want %d", webErr.GotStatus, http.StatusNotFound)
	}
	if webErr.Method != "GET" {
		t.Errorf("Method = %s, want GET", webErr.Method)
	}
	if webErr.ResponseBody != "resource not found" {
		t.Errorf("ResponseBody = %s, want 'resource not found'", webErr.ResponseBody)
	}
	if webErr.URL == nil {
		t.Fatal("URL is nil")
	}
	if !strings.Contains(webErr.Error(), "resource not found") {
		t.Errorf("Error() string does not contain response body: %s", webErr.Error())
	}
}
