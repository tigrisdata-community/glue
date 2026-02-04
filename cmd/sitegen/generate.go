package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GeneratePages(cfg *Config, pages []*Page) error {
	for _, page := range pages {
		// Calculate output paths
		relPath, err := filepath.Rel(cfg.ContentDir, page.InputPath)
		if err != nil {
			return err
		}

		htmlPath := filepath.Join(cfg.OutputDir, strings.TrimSuffix(relPath, ".md")+".html")
		mdPath := filepath.Join(cfg.OutputDir, relPath)

		// Create output directory
		if err := os.MkdirAll(filepath.Dir(htmlPath), 0755); err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}

		// Render HTML
		baseTempl := Base(page.Frontmatter.Title, PageView(page.Frontmatter.Title, page.Frontmatter.Body))
		var buf bytes.Buffer
		if err := baseTempl.Render(context.Background(), &buf); err != nil {
			return fmt.Errorf("rendering template: %w", err)
		}

		if err := os.WriteFile(htmlPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("writing HTML: %w", err)
		}

		// Copy original markdown
		if err := os.MkdirAll(filepath.Dir(mdPath), 0755); err != nil {
			return err
		}
		content, err := os.ReadFile(page.InputPath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(mdPath, content, 0644); err != nil {
			return err
		}
	}

	return nil
}
