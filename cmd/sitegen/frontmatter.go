package main

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

type Frontmatter struct {
	Title       string
	Description string
	Body        string
}

func ParseFrontmatter(content []byte) (*Frontmatter, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{
				Formats: []frontmatter.Format{frontmatter.YAML},
				Mode:    frontmatter.SetMetadata,
			},
		),
	)

	context := parser.NewContext()
	reader := text.NewReader(content)
	md.Parser().Parse(reader, parser.WithContext(context))

	// Get metadata from parser context
	data := frontmatter.Get(context)
	if data == nil {
		return nil, fmt.Errorf("frontmatter not found or empty")
	}

	var meta map[string]string
	if err := data.Decode(&meta); err != nil {
		return nil, fmt.Errorf("decoding frontmatter: %w", err)
	}

	if len(meta) == 0 {
		return nil, fmt.Errorf("frontmatter not found or empty")
	}

	title, ok := meta["title"]
	if !ok || title == "" {
		return nil, fmt.Errorf("title is required in frontmatter")
	}

	description, ok := meta["description"]
	if !ok || description == "" {
		return nil, fmt.Errorf("description is required in frontmatter")
	}

	var body bytes.Buffer
	if err := md.Convert(content, &body); err != nil {
		return nil, fmt.Errorf("converting markdown: %w", err)
	}

	return &Frontmatter{
		Title:       title,
		Description: description,
		Body:        body.String(),
	}, nil
}
