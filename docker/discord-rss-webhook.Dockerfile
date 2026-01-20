FROM golang:1.25 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache \
  GOBIN=/app/bin \
  go install -ldflags="-X github.com/tigrisdata-community/glue=$(git describe --tags --always --dirty)" \
  ./cmd/discord-rss-webhook

FROM debian:bookworm AS runtime
WORKDIR /app

RUN apt-get update && apt-get install -y \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/* \
  && cp /etc/ssl/certs/ca-certificates.crt .

COPY --from=build /app/bin/discord-rss-webhook /app/bin/discord-rss-webhook
CMD ["/app/bin/discord-rss-webhook"]

LABEL org.opencontainers.image.source="https://github.com/tigrisdata-community/glue"

