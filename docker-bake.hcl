variable "VERSION" { default = "devel" }

group "default" {
  targets = [
    "discord-rss-webhook",
  ]
}

target "discord-rss-webhook" {
  context = "."
  dockerfile = "./docker/discord-rss-webhook.Dockerfile"
  platforms = [
    "linux/amd64",
    "linux/arm64",
  ]
  pull = true
  tags = [
    "ghcr.io/tigrisdata-community/glue/discord-rss-webhook:${VERSION}",
    "ghcr.io/tigrisdata-community/glue/discord-rss-webhook:latest",
  ]
}