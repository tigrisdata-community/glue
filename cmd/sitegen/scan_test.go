package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanContent(t *testing.T) {
	// Create temporary content directory
	tmpDir, err := os.MkdirTemp("", "sitegen-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := map[string]string{
		"index.md": `---
title: "Index"
description: "Main index"
---
Index content`,
		"guide/index.md": `---
title: "Guide"
description: "Guide index"
---
Guide content`,
		"guide/setup.md": `---
title: "Setup"
description: "Setup guide"
---
Setup content`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	pages, err := ScanContent(tmpDir)
	if err != nil {
		t.Fatalf("ScanContent() error = %v", err)
	}

	// Should find 3 files
	if len(pages) != 3 {
		t.Errorf("ScanContent() found %d files, want 3", len(pages))
	}

	// Check that index.md files are identified
	var indexCount int
	for _, p := range pages {
		if filepath.Base(p.InputPath) == "index.md" {
			indexCount++
		}
	}
	if indexCount != 2 {
		t.Errorf("ScanContent() found %d index.md files, want 2", indexCount)
	}
}
