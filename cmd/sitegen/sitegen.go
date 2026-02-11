package main

import (
	"fmt"
	"log"
	"os"
)

func Generate(cfg *Config, quiet bool) error {
	// Create output directory
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Scan content directory
	if !quiet {
		log.Printf("Scanning %s...", cfg.ContentDir)
	}
	pages, err := ScanContent(cfg.ContentDir)
	if err != nil {
		return fmt.Errorf("scanning content: %w", err)
	}

	if !quiet {
		log.Printf("Found %d markdown files", len(pages))
	}

	// Parse frontmatter for each page
	for i, page := range pages {
		if !quiet {
			log.Printf("Parsing %s...", page.InputPath)
		}
		content, err := os.ReadFile(page.InputPath)
		if err != nil {
			return fmt.Errorf("reading %s: %w", page.InputPath, err)
		}

		fm, err := ParseFrontmatter(content)
		if err != nil {
			return fmt.Errorf("parsing frontmatter in %s: %w", page.InputPath, err)
		}
		pages[i].Frontmatter = fm
	}

	// Generate HTML pages and copy markdown
	if !quiet {
		log.Printf("Generating HTML...")
	}
	if err := GeneratePages(cfg, pages); err != nil {
		return fmt.Errorf("generating pages: %w", err)
	}

	// Generate llms.txt
	if !quiet {
		log.Printf("Generating llms.txt...")
	}
	if err := GenerateLLMsTxt(cfg, pages); err != nil {
		return fmt.Errorf("generating llms.txt: %w", err)
	}

	return nil
}
