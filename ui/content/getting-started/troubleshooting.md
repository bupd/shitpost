---
title: "Troubleshooting"
description: "Common startup, auth, media, and platform posting failures."
weight: 80
---

## Bot fails on startup

Check `BOT_TOKEN` first. If it is missing, the process exits immediately with:

```text
BOT_TOKEN environment variable is required
```

If the token is wrong, Telegram bot creation fails. Regenerate or copy the token again from BotFather.

## Bot starts but ignores you

If `AUTHORIZED_TELEGRAM_USERS` is set, your Telegram username or numeric user ID must be in the comma-separated list.

```dotenv
AUTHORIZED_TELEGRAM_USERS=your_username,123456789
```

Remove the leading `@` or keep it; `shitpost` normalizes both forms.

## Posts do not appear

Run the doctor script:

```sh
task doctor
```

Confirm the target platform variables are present. Then run dry-run mode and inspect the command preview:

```sh
task up:dry-run
```

If the command preview is wrong, fix `CROSSPOST_FLAGS`. If the command preview is right but live posting fails, inspect the `crosspost` stderr returned in Telegram.

## X returns 401

Confirm you are using the correct X credential path.

For the emusks-backed strategy, set only the cookie value:

```dotenv
AUTH_TOKEN=your-auth-token-cookie-value
```

For official API fallback, verify `TWITTER_ACCESS_TOKEN_KEY` is an OAuth 1.0a access token. The doctor script warns if the value does not look like a typical OAuth token.

## Media does not upload

Check that `downloads/` exists and is writable by the container.

```sh
mkdir -p downloads
docker compose up --build
```

Images are attached through `crosspost`. Videos and non-image documents currently post caption text only.

## Logs are too long for Telegram

Telegram has a message length limit. `shitpost` chunks long replies so you still receive the command result and logs.

## Reset safely

Stop and remove the container without deleting your `.env` or downloaded media:

```sh
docker rm -f shitpost
docker run -d --name shitpost \
  --env-file .env \
  -v ./downloads:/app/downloads \
  ghcr.io/bupd/shitpost:latest
```
