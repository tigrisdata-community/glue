package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
)

func GeneratePages(cfg *Config, pages []*Page) error {
	// Group pages by directory for index page lists
	dirPages := make(map[string][]*Page)
	for _, page := range pages {
		dir := filepath.Dir(page.URLPath)
		if dir == "." {
			dir = ""
		}
		if !page.IsIndex {
			dirPages[dir] = append(dirPages[dir], page)
		}
	}

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

		// Build page list for index files
		var content templ.Component
		if page.IsIndex {
			dir := filepath.Dir(page.URLPath)
			if dir == "." {
				dir = ""
			}
			childPages := dirPages[dir]

			var listBuilder strings.Builder
			listBuilder.WriteString("<ul>")
			for _, cp := range childPages {
				href := strings.TrimPrefix(cp.URLPath, "./")
				desc := cleanDescription(cp.Frontmatter.Description)
				listBuilder.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a>: %s</li>", href, cp.Frontmatter.Title, desc))
			}
			listBuilder.WriteString("</ul>")
			pageListHTML := listBuilder.String()

			content = PageIndex(page.Frontmatter.Title, page.Frontmatter.Body, pageListHTML)
		} else {
			content = PageView(page.Frontmatter.Title, page.Frontmatter.Body)
		}

		// Render HTML
		baseTempl := Base(page.Frontmatter.Title, content)
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
		contentBytes, err := os.ReadFile(page.InputPath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(mdPath, contentBytes, 0644); err != nil {
			return err
		}
	}

	return nil
}
