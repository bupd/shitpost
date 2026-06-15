# ----------------------------------------
# STAGE 1 — BUILD GO BINARY
# ----------------------------------------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bot main.go


# ----------------------------------------
# STAGE 2 — RUNTIME WITH BUN & CROSSPOST
# ----------------------------------------
FROM oven/bun:1.3.10

WORKDIR /app

ARG CROSSPOST_REPO=https://github.com/bupd/crosspost.git
ARG CROSSPOST_REF=main

RUN apt-get update && \
    apt-get install -y --no-install-recommends git ca-certificates && \
    git clone --depth=1 --branch "$CROSSPOST_REF" "$CROSSPOST_REPO" /opt/crosspost && \
    cd /opt/crosspost && \
    bun install --frozen-lockfile && \
    bun run build && \
    printf '#!/usr/bin/env sh\nexec bun /opt/crosspost/dist/bin.js "$@"\n' > /usr/local/bin/crosspost && \
    chmod +x /usr/local/bin/crosspost && \
    apt-get purge -y --auto-remove git && \
    rm -rf /var/lib/apt/lists/*

# Mastodon rejects Bun FormData Blob uploads without a filename as blank files.
RUN sed -i 's/data.append("file", new Blob(\[image.data\], { type }));/data.append("file", new Blob([image.data], { type }), "image.jpg");/' \
    /opt/crosspost/dist/strategies/mastodon.js

COPY --from=builder /app/bot /app/bot

RUN mkdir -p /app/downloads

CMD ["./bot"]
