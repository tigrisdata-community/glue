package main

import (
	"fmt"
	"os"
)

func Generate(cfg *Config, quiet bool) error {
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// TODO: Scan content directory
	// TODO: Parse frontmatter
	// TODO: Generate HTML
	// TODO: Generate llms.txt

	return nil
}
