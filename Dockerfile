# ----------------------------------------
# STAGE 1 — BUILD GO BINARY
# ----------------------------------------
FROM golang:1.25 AS builder

WORKDIR /app

# Copy dependencies first (better cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy app source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o bot main.go


# ----------------------------------------
# STAGE 2 — RUNTIME WITH NODE & CROSSPOST
# ----------------------------------------
FROM debian:stable-slim

WORKDIR /app

# Install Node + npm
RUN apt-get update && apt-get install -y curl ca-certificates \
    && curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# Install the scoped package globally
RUN npm install -g @humanwhocodes/crosspost

# Copy Go binary from builder stage
COPY --from=builder /app/bot /app/bot

# Create downloads dir for media
RUN mkdir -p /app/downloads

# Run the Go bot
CMD ["./bot"]
