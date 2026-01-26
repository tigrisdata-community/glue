#!/bin/sh
# Build and deploy discord-rss-webhook binary to target machine
# Usage: ./build.sh [user@]host

set -e

TARGET="${1:?Usage: $0 [user@]host}"

echo "Building discord-rss-webhook for linux/amd64..."
GOOS=linux GOARCH=amd64 go build -o discord-rss-webhook ./cmd/discord-rss-webhook

echo "Installing to ${TARGET}:/home/exedev/bin/..."
ssh "${TARGET}" "mkdir -p /home/exedev/bin"
scp discord-rss-webhook "${TARGET}:/home/exedev/bin/"
ssh "${TARGET}" "chmod +x /home/exedev/bin/discord-rss-webhook"

rm -f discord-rss-webhook
echo "Done!"
