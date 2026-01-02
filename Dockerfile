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
FROM oven/bun:alpine

WORKDIR /app

RUN bun install -g @humanwhocodes/crosspost

COPY --from=builder /app/bot /app/bot

RUN mkdir -p /app/downloads

CMD ["./bot"]
