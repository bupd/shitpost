# shitpost

[![CI](https://github.com/bupd/shitpost/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/bupd/shitpost/actions/workflows/ci.yml)

Telegram-powered social media crossposting tool. Post once to Telegram, publish everywhere.

Accept text, photos, videos, and documents via a [Telegram bot](https://t.me/shitpost_engine_bot), then publish them across all your social platforms simultaneously.

A lightweight, self-hosted alternative to Postiz, Buffer, and Hootsuite for developers who want full control over their social media automation.

## Motivation

- I am fed up with posting and managing different platforms.
- I have multiple thoughts popping up in my head which I would like to convey to people
- But I lack the patience to go through multiple different apps just to say something...
- And to update the content based on the app is more hated.

## Features

- Post to multiple social networks simultaneously from Telegram
- Supports text, images, videos, and documents
- Alt-text support for accessible image posts
- Self-hosted and privacy-focused
- Multi-arch Docker images (amd64, arm64)
- Lightweight Go binary with minimal dependencies

## Supported Platforms

| Platform | Text | Images | Videos |
|----------|------|--------|--------|
| Twitter/X | Yes | Yes | Yes |
| Bluesky | Yes | Yes | Yes |
| Mastodon | Yes | Yes | Yes |
| LinkedIn | Yes | Yes | No |
| Discord | Yes | Yes | Yes |
| Telegram | Yes | Yes | Yes |
| Slack | Yes | Yes | No |
| Dev.to | Yes | No | No |
| Nostr | Yes | No | No |

## Why shitpost?

| Feature | shitpost | Postiz | Buffer |
|---------|----------|--------|--------|
| Self-hosted | Yes | Yes | No |
| Free | Yes | Freemium | Freemium |
| Telegram interface | Yes | No | No |
| No web UI required | Yes | No | No |
| Privacy-focused | Yes | Partial | No |
| Open source | Yes | Yes | No |

## Docker Images

Pre-built multi-arch images (amd64, arm64) are available:

```sh
# Primary (Harbor)
docker pull registry.goharbor.io/bupd/shitpost:latest

# Alternative (GitHub Container Registry)
docker pull ghcr.io/bupd/shitpost:latest
```

### Image Tags

| Tag | Description |
|-----|-------------|
| `latest` | Most recent build from main branch. May include untested changes. |
| `v1.0.0` | Specific release version. Stable and tested. Recommended for production. |
| `v1.0` | Latest patch release in v1.0.x series. |
| `v1` | Latest minor release in v1.x.x series. |

For production, use a versioned tag (e.g., `v1.0.0`) to avoid unexpected updates.

## Quickstart

### Option 1: Run pre-built image (recommended)

1. Create `.env` file:
   ```sh
   curl -o .env https://raw.githubusercontent.com/bupd/shitpost/main/.env.example
   ```
   Edit `.env` and set your tokens.

2. Run the container:
   ```sh
   docker run -d --name shitpost \
     --env-file .env \
     -v ./downloads:/app/downloads \
     registry.goharbor.io/bupd/shitpost:latest
   ```

3. Check logs:
   ```sh
   docker logs -f shitpost
   ```

### Option 2: Build from source

1. Clone and enter repo:
   ```sh
   git clone https://github.com/bupd/shitpost.git
   cd shitpost
   ```

2. Create `.env` file from template and set `BOT_TOKEN`:
   ```sh
   cp .env.example .env
   ```

3. Build & run with Docker Compose:
   ```sh
   docker compose up --build
   ```

4. Confirm the bot is running by checking logs for:
   `Authorized as @<your_bot_username>`

5. Send messages or media to your bot in Telegram. The bot will post using crosspost and reply with logs.

## Requirements

- Docker/Podman (recommended) or Go 1.25+
- Telegram bot token (create via [@BotFather](https://t.me/BotFather))
- API keys for target platforms (Twitter, Bluesky, Mastodon, etc.)

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `BOT_TOKEN` | Yes | Telegram bot token from BotFather |
| `TWITTER_API_CONSUMER_KEY` | No | Twitter API consumer key |
| `TWITTER_API_CONSUMER_SECRET` | No | Twitter API consumer secret |
| `TWITTER_ACCESS_TOKEN_KEY` | No | Twitter access token |
| `TWITTER_ACCESS_TOKEN_SECRET` | No | Twitter access token secret |
| `BLUESKY_HOST` | No | Bluesky host (e.g., bsky.social) |
| `BLUESKY_IDENTIFIER` | No | Bluesky handle or email |
| `BLUESKY_PASSWORD` | No | Bluesky app password |
| `MASTODON_HOST` | No | Mastodon instance URL |
| `MASTODON_ACCESS_TOKEN` | No | Mastodon access token |

crosspost itself may require additional environment variables (API keys, tokens for target platforms). For details about those envs and how to obtain them, consult the [crosspost documentation](https://github.com/humanwhocodes/crosspost).

## Usage

### Text Posts
Send any text message to your bot. It will be posted to all configured platforms.

### Media Posts
Send images, videos, or documents with an optional caption.

### Alt Text
Add alt text for accessibility by ending your caption with `alt:`:
```
Check out this sunset!
alt: Orange and purple sunset over mountains
```

## Running Locally (without Docker)

1. Install Go 1.25+ and crosspost CLI:
   ```sh
   npm install -g @humanwhocodes/crosspost
   ```

2. Build and run:
   ```sh
   go build -o bot main.go
   export BOT_TOKEN=your_token_here
   ./bot
   ```

## Architecture

```
Telegram → shitpost bot → crosspost CLI → Social platforms
```

Built as a lightweight wrapper around [humanwhocodes/crosspost](https://github.com/humanwhocodes/crosspost).

## Security

- Self-hosted: your data stays on your server
- No third-party analytics or tracking
- Media files stored locally (configure volume mounts)
- Consider implementing allowlists for public-facing bots

## Troubleshooting

### Bot fails on startup
- Verify `BOT_TOKEN` is set correctly
- Check logs: `docker compose logs -f`

### Posts not appearing
- Verify platform API keys are configured
- Check crosspost output in bot logs

### Media not uploading
- Ensure `./downloads` directory exists and is writable
- Check disk space

## Contributing

Contributions welcome! Please open an issue or PR.

## License

See [LICENSE](LICENSE) in the repository.

## Related Projects

- [crosspost](https://github.com/humanwhocodes/crosspost) - CLI tool powering the cross-posting
- [Postiz](https://github.com/gitroomhq/postiz-app) - Full-featured social media scheduler
- [Buffer](https://buffer.com) - Commercial social media management

## GitHub Topics

Add these topics to your repository for better discoverability:
`crosspost`, `social-media`, `telegram-bot`, `twitter`, `bluesky`, `mastodon`, `automation`, `self-hosted`, `golang`, `docker`
