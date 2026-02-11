package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
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

	// Create test content structure
	files := map[string]string{
		"index.md": `---
title: "Home"
description: "Welcome to the site"
---
# Welcome

This is the home page.`,
		"about.md": `---
title: "About"
description: "About this site"
---
# About

Information about the site.`,
		"docs/index.md": `---
title: "Documentation"
description: "Main docs index"
---
# Docs

Documentation index.`,
	}

	for path, content := range files {
		fullPath := filepath.Join(contentDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	cfg := &Config{
		ContentDir: contentDir,
		OutputDir:  outputDir,
		Preamble:   "# My Site\n\nDocumentation site.\n",
	}

	// Run generation
	if err := Generate(cfg, true); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check HTML files were generated
	htmlFiles := []string{"index.html", "about.html", "docs/index.html"}
	for _, f := range htmlFiles {
		path := filepath.Join(outputDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("HTML file not generated: %s", f)
		}
	}

	// Check markdown files were copied
	mdFiles := []string{"index.md", "about.md", "docs/index.md"}
	for _, f := range mdFiles {
		path := filepath.Join(outputDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Markdown file not copied: %s", f)
		}
	}

	// Check llms.txt
	llmsPath := filepath.Join(outputDir, "llms.txt")
	llmsContent, err := os.ReadFile(llmsPath)
	if err != nil {
		t.Fatalf("Reading llms.txt: %v", err)
	}
	llmsStr := string(llmsContent)
	if !strings.Contains(llmsStr, "# My Site") {
		t.Errorf("llms.txt missing preamble")
	}
	// index.md files are skipped in llms.txt
	if strings.Contains(llmsStr, "[Home](./index.md)") {
		t.Errorf("llms.txt should not contain index.md links")
	}
	if !strings.Contains(llmsStr, "[About](./about.md)") {
		t.Errorf("llms.txt missing About link")
	}
}
