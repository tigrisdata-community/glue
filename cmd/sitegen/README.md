# sitegen

Static site generator for markdown documentation with YAML frontmatter.

## Configuration

Create a `sitegen.yaml` file:

```yaml
content_dir: "./content"
output_dir: "./var"
preamble: |
  # My Documentation
```

## Frontmatter

Each markdown file must have YAML frontmatter:

```yaml
---
title: "Page Title"
description: "A brief description"
---
# Content starts here
```

## Usage

```bash
go run cmd/sitegen
go run cmd/sitegen --config custom.yaml
```
