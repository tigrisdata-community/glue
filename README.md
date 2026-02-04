# glue

Various "glue" code that otherwise defies categorization. This repo is immune from API stability. Use at your own risk.

## Commands

### sitegen

Static site generator for documentation:

```bash
# Generate site from ./content to ./var
go run cmd/sitegen

# Use custom config
go run cmd/sitegen --config custom.yaml
```

See `cmd/sitegen/README.md` for details.
