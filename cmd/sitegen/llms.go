package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func cleanDescription(desc string) string {
	// Remove newlines
	desc = strings.ReplaceAll(desc, "\n", " ")
	// Collapse multiple spaces to single space
	spaceRegex := regexp.MustCompile(`\s+`)
	desc = spaceRegex.ReplaceAllString(desc, " ")
	// Trim leading/trailing whitespace
	desc = strings.TrimSpace(desc)
	return desc
}

func GenerateLLMsTxt(cfg *Config, pages []*Page) error {
	var sb strings.Builder

	// Write preamble
	if cfg.Preamble != "" {
		sb.WriteString(cfg.Preamble)
		if !strings.HasSuffix(cfg.Preamble, "\n") {
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Write entries
	for _, page := range pages {
		title := page.Frontmatter.Title
		link := page.URLPath
		desc := cleanDescription(page.Frontmatter.Description)

		sb.WriteString(fmt.Sprintf("[%s](%s): %s\n", title, link, desc))
	}

	// Write to file
	path := filepath.Join(cfg.OutputDir, "llms.txt")
	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("writing llms.txt: %w", err)
	}

	return nil
}
