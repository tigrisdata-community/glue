# Popola

Popola generates Tigris tutorials autonomously. It reads a topic from stdin, consults documentation, and instructs Claude to write a hands-on tutorial.

## Input

Pipe JSON to stdin with two fields:

```json
{
  "topic": "Migrating data from Hetzner object storage to Tigris",
  "relevantDocs": "https://www.tigrisdata.com/docs/migration/\n* https://docs.hetzner.com/storage/object-storage/getting-started/using-s3-api-tools/"
}
```

Both fields are required.

## Output

The agent writes Markdown tutorials to `./var` in a sensible location, then prints the final file path.

## How it works

Popola hydrates a prompt template with your input, then launches a Claude Code session. The agent can read the web and query the Tigris community Discord. It writes the tutorial, revises it for clarity and style, then saves it to disk.

Tool usage and events stream to stdout as JSON lines.

## Flags

| Flag                    | Description                  | Default                  |
| ----------------------- | ---------------------------- | ------------------------ |
| `-anthropic-auth-token` | Anthropic API token          | `hunter2`                |
| `-anthropic-base-url`   | Anthropic API base URL       | `http://localhost:11434` |
| `-anthropic-model`      | Model to use                 | `glm-4.7-flash:latest`   |
| `-zhipu-api-key`        | Zhipu API key for web reader | (empty)                  |

## Example

```bash
echo '{"topic":"Migrating from S3 to Tigris","relevantDocs":"* https://www.tigrisdata.com/docs/migration/"}' | \
  go run ./cmd/popola
```
