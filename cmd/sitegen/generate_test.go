package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratePages(t *testing.T) {
	// Create temporary directories
	contentDir, err := os.MkdirTemp("", "sitegen-content-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(contentDir)

	outputDir, err := os.MkdirTemp("", "sitegen-output-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(outputDir)

	// Create test markdown
	mdContent := `---
title: "Test Page"
description: "A test"
---
# Hello World`

	mdPath := filepath.Join(contentDir, "test.md")
	if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Scan and parse
	pages, err := ScanContent(contentDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range pages {
		content, err := os.ReadFile(p.InputPath)
		if err != nil {
			t.Fatal(err)
		}
		p.Frontmatter, err = ParseFrontmatter(content)
		if err != nil {
			t.Fatal(err)
		}
	}

	cfg := &Config{ContentDir: contentDir, OutputDir: outputDir}

	// Generate
	if err := GeneratePages(cfg, pages); err != nil {
		t.Fatalf("GeneratePages() error = %v", err)
	}

	// Check HTML output exists
	htmlPath := filepath.Join(outputDir, "test.html")
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		t.Errorf("GeneratePages() did not create %s", htmlPath)
	}

	// Check HTML contains expected content
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		t.Fatal(err)
	}
	htmlStr := string(htmlContent)
	if !strings.Contains(htmlStr, "Test Page") {
		t.Errorf("HTML does not contain title 'Test Page'")
	}
	if !strings.Contains(htmlStr, "Hello World") {
		t.Errorf("HTML does not contain 'Hello World'")
	}

	// Check markdown was copied
	mdOutputPath := filepath.Join(outputDir, "test.md")
	if _, err := os.Stat(mdOutputPath); os.IsNotExist(err) {
		t.Errorf("GeneratePages() did not copy %s", mdOutputPath)
	}
}
