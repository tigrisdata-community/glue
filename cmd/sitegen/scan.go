package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
	InputPath   string // Full path to source .md file
	OutputPath  string // Full path to output .html file
	OutputMD    string // Full path to copied .md file
	URLPath     string // Relative path for linking (e.g., "./guide/setup.md")
	IsIndex     bool   // True if filename is index.md
	Frontmatter *Frontmatter
}

func ScanContent(contentDir string) ([]*Page, error) {
	var pages []*Page

	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		relPath, err := filepath.Rel(contentDir, path)
		if err != nil {
			return err
		}

		pages = append(pages, &Page{
			InputPath: path,
			URLPath:   "./" + relPath,
			IsIndex:   filepath.Base(path) == "index.md",
		})

		return nil
	})

	return pages, err
}
