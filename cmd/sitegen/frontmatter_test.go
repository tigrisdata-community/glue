package main

import (
	"strings"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantTitle   string
		wantDesc    string
		wantBody    string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid frontmatter",
			content: `---
title: "Hello World"
description: "A test page"
---

# Content here`,
			wantTitle: "Hello World",
			wantDesc:  "A test page",
			wantBody:  "<h1>Content here</h1>\n",
			wantErr:   false,
		},
		{
			name:        "missing frontmatter",
			content:     `# No frontmatter here`,
			wantErr:     true,
			errContains: "frontmatter not found",
		},
		{
			name: "missing title",
			content: `---
description: "No title"
---

Content`,
			wantErr:     true,
			errContains: "title is required",
		},
		{
			name: "missing description",
			content: `---
title: "No description"
---

Content`,
			wantErr:     true,
			errContains: "description is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFrontmatter([]byte(tt.content))
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseFrontmatter() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ParseFrontmatter() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseFrontmatter() unexpected error: %v", err)
				return
			}
			if got.Title != tt.wantTitle {
				t.Errorf("ParseFrontmatter() Title = %q, want %q", got.Title, tt.wantTitle)
			}
			if got.Description != tt.wantDesc {
				t.Errorf("ParseFrontmatter() Description = %q, want %q", got.Description, tt.wantDesc)
			}
			if got.Body != tt.wantBody {
				t.Errorf("ParseFrontmatter() Body = %q, want %q", got.Body, tt.wantBody)
			}
		})
	}
}
