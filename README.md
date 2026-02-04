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

### popola

Popola serves as the canonical implementation of the Omnlana "agent protocol", enabling
consistent agent interactions across AI platforms. It reads a JSON input from stdin with
`topic` and `relevantDocs` fields, then launches a Claude Code session with MCP tools for
web reading and Tigris Discord integration.

See `cmd/popola/README.md` for details.
