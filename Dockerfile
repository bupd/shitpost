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
FROM oven/bun:1.3.10-alpine

WORKDIR /app

RUN bun install -g @humanwhocodes/crosspost

# Mastodon rejects Bun FormData Blob uploads without a filename as blank files.
RUN sed -i 's/data.append("file", new Blob(\[image.data\], { type }));/data.append("file", new Blob([image.data], { type }), "image.jpg");/' \
    /root/.bun/install/global/node_modules/@humanwhocodes/crosspost/dist/strategies/mastodon.js

COPY --from=builder /app/bot /app/bot

RUN mkdir -p /app/downloads

CMD ["./bot"]
