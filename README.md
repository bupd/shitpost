# shitpost

Simple crossposting bridge that accepts text, photos, videos, and documents via a [Telegram bot](https://t.me/shitpost_engine_bot), runs the [crosspost](https://github.com/humanwhocodes/crosspost) CLI to publish them across social platforms.

## Quick summary
- Listens for posts using Telegram Bot.
- Calls the `crosspost` CLI to publish (text-only or with image + alt text).
- Built as a lightweight abstraction over [humanwhocodes/crosspost](https://github.com/humanwhocodes/crosspost) (see that repo for crosspost-specific configuration).


## Motivation
- I am fed up with posting and managing different platforms.
- I have multiple thoughts popping up in my head which I would like to convey to people
- But I lack the patience to go through multiple different apps just to say something...
- And to update the content based on the app is more hated.

  
## Requirements
- Podman/Docker & Compose (recommended) or Go 1.25+ to build/run locally.
- A Telegram bot token (create via @BotFather).

## Quickstart
1. Clone and enter repo:
```sh
   git clone https://github.com/bupd/shitpost.git
   cd shitpost
```

2. Create `.env` file from template `.env.example` and set `BOT_TOKEN`:
   ```sh
   cp .env.example .env
   ```
   # Edit .env and set:
   ```
   BOT_TOKEN=123456:ABC-DEF...
   ```

3. Build & run with Docker Compose:
   ```
   docker compose up --build
   ```

4. Confirm the bot is running by checking logs for:
   Authorized as @<your_bot_username>

5. Send messages or media to your bot in Telegram. The bot will post using crosspost and reply with logs / the same file.


## Environment variables

- BOT_TOKEN (required) — Telegram bot token from BotFather.

crosspost itself may require additional environment variables (API keys, tokens for target platforms). For details about those envs and how to obtain them, consult:
https://github.com/humanwhocodes/crosspost

Use .env (or docker-compose env_file) to supply values.


## Telegram usage / caption rules

- Text messages: posted as text via crosspost.
- Images/videos/documents: downloaded and posted via crosspost.
- Alt-text parsing: if the last line of the caption starts with `alt:` (case-insensitive), that line is removed from the caption and used as the image alt text.
  Example:
  ```
  Here’s the pic
  alt: A smiling cat on a red blanket
  ```


## Running locally (without Docker)

1. Ensure Go 1.25+ is installed.
2. Install dependencies:
   ```
   go mod download
   ```
3. Build:
   ```
   CGO_ENABLED=0 GOOS=linux go build -o bot main.go
   ```
4. Export BOT_TOKEN and run:
   ```
   export BOT_TOKEN=123456:ABC-DEF...
   ./bot
   ```

Note: When running locally, make sure the `crosspost` CLI is available in your PATH (install it with npm: npm install -g @humanwhocodes/crosspost) or adjust PATH accordingly.


## Persistence & volumes
- Media downloaded from Telegram are stored at ./downloads in the repo root (mounted to /app/downloads in the container). Keep this folder secure or change the mapping if needed.


## Security & privacy notes
- The bot downloads user media to local storage — protect the host and mounted volumes.
- crosspost stdout/stderr are returned to the chat; ensure crosspost does not leak secrets or sensitive tokens in logs.
- Consider implementing user allowlists if the bot will be public-facing (not implemented by default).



## Troubleshooting

- Bot fails on startup:
  - Ensure BOT_TOKEN is set and valid.
  - Check container logs: docker compose logs -f

- crosspost not found:
  - Dockerfile installs crosspost globally; if running locally, install it with:
    npm install -g @humanwhocodes/crosspost
  - Or ensure crosspost binary is on PATH.

- Files not saved:
  - Ensure ./downloads exists and is writeable by the container/user.
  - The compose file maps ./downloads to /app/downloads; permissions on host may need adjusting.


## Development notes & roadmap (short)
- Add allowlist/denylist for users or groups.
- Add config options to control whether the posted file is returned.
- Add graceful shutdown handling.
- Add tests and CI.
- Improve logging and error reporting.

## License
See LICENSE in the repository.

References
- crosspost: https://github.com/humanwhocodes/crosspost

