package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateLLMsTxt(t *testing.T) {
	// Create temporary output directory
	outputDir, err := os.MkdirTemp("", "sitegen-llms-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(outputDir)

	cfg := &Config{
		OutputDir: outputDir,
		Preamble:  "# Test Docs\n\nWelcome to the docs.\n",
	}

	pages := []*Page{
		{
			URLPath: "./index.md",
			Frontmatter: &Frontmatter{
				Title:       "Main Page",
				Description: "The main index page",
			},
		},
		{
			URLPath: "./guide/setup.md",
			Frontmatter: &Frontmatter{
				Title:       "Setup Guide",
				Description: "How to get started",
			},
		},
		{
			URLPath: "./api/reference.md",
			Frontmatter: &Frontmatter{
				Title:       "API Reference",
				Description: "API docs with   extra    spaces\nand newlines",
			},
		},
	}

	if err := GenerateLLMsTxt(cfg, pages); err != nil {
		t.Fatalf("GenerateLLMsTxt() error = %v", err)
	}

	// Check file exists
	llmsPath := filepath.Join(outputDir, "llms.txt")
	content, err := os.ReadFile(llmsPath)
	if err != nil {
		t.Fatalf("Reading llms.txt: %v", err)
	}

	contentStr := string(content)

	// Check preamble
	if !strings.Contains(contentStr, "# Test Docs") {
		t.Errorf("llms.txt missing preamble")
	}

	// Check link format
	if !strings.Contains(contentStr, "[Main Page](./index.md): The main index page") {
		t.Errorf("llms.txt missing correct link format for Main Page")
	}

	// Check description cleanup (extra spaces and newlines removed)
	if strings.Contains(contentStr, "extra    spaces") {
		t.Errorf("llms.txt description not cleaned of extra spaces")
	}
}

func TestCleanDescription(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal text",
			input:    "Normal description here",
			expected: "Normal description here",
		},
		{
			name:     "extra spaces",
			input:    "Text    with    extra    spaces",
			expected: "Text with extra spaces",
		},
		{
			name:     "newlines",
			input:    "Text\nwith\nnewlines",
			expected: "Text with newlines",
		},
		{
			name:     "mixed",
			input:    "Text\nwith   both\n   issues",
			expected: "Text with both issues",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanDescription(tt.input)
			if got != tt.expected {
				t.Errorf("cleanDescription() = %q, want %q", got, tt.expected)
			}
		})
	}
}
