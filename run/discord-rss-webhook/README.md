# Discord RSS Webhook - Systemd Setup

## Installation

1. Copy the units to systemd:

   ```sh
   sudo cp discord-rss-webhook.{service,timer} /etc/systemd/system/
   ```

2. Create the environment file at `/home/exedev/discord-rss-webhook.env`:

   ```sh
   STORE_BUCKET=your_bucket
   AWS_ACCESS_KEY_ID=your_key
   AWS_SECRET_ACCESS_KEY=your_secret
   AWS_ENDPOINT_URL_S3=your_endpoint
   AWS_ENDPOINT_URL_IAM=your_endpoint
   AWS_REGION=your_region
   DISCORD_WEBHOOK_URL=your_webhook_url
   ```

3. Enable and start the timer:
   ```sh
   sudo systemctl daemon-reload
   sudo systemctl enable --now discord-rss-webhook.timer
   ```

## Usage

- Check timer status: `systemctl status discord-rss-webhook.timer`
- View next run times: `systemctl list-timers discord-rss-webhook.timer`
- Run manually: `sudo systemctl start discord-rss-webhook.service`
- View logs: `journalctl -u discord-rss-webhook.service -f`
