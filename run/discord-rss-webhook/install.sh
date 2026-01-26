#!/bin/sh
# Install systemd service and timer to target machine
# Usage: ./install.sh [user@]host

set -e

TARGET="${1:?Usage: $0 [user@]host}"

echo "Installing discord-rss-webhook to $TARGET..."

# Copy systemd units
scp run/discord-rss-webhook/discord-rss-webhook.service run/discord-rss-webhook/discord-rss-webhook.timer "${TARGET}:/tmp/"

# Install and enable on remote
ssh "${TARGET}" <<'EOF'
set -e
sudo mv /tmp/discord-rss-webhook.service /tmp/discord-rss-webhook.timer /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now discord-rss-webhook.timer
echo "Installation complete. Timer status:"
systemctl status discord-rss-webhook.timer --no-pager
EOF

echo "Done! Don't forget to create /home/exedev/discord-rss-webhook.env on the target."
