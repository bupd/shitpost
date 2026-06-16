---
title: "Getting started"
description: "Run a private Telegram bot that crossposts text and media to your social accounts."
weight: 10
---

`shitpost` is a small Go service that listens to your Telegram bot, downloads incoming media when needed, and calls the `crosspost` CLI with the platform flags you choose.

The happy path is:

1. Create a Telegram bot with BotFather.
2. Create a `.env` file from `.env.example`.
3. Add your Telegram token and the platform credentials you need.
4. Start in dry-run mode.
5. Send a message to your bot and inspect the command preview.
6. Disable dry-run when the output looks right.

```sh
curl -o .env https://raw.githubusercontent.com/bupd/shitpost/main/.env.example
docker run -d --name shitpost \
  --env-file .env \
  -v ./downloads:/app/downloads \
  ghcr.io/bupd/shitpost:latest
```

## What it does

- Accepts Telegram text messages and posts them through `crosspost`.
- Accepts Telegram photos, downloads the largest image, and passes it to `crosspost` with optional alt text.
- Accepts video and document uploads, stores them locally, and posts caption text when the installed `crosspost` media path cannot attach them.
- Replies back in Telegram with the `crosspost` logs or the dry-run command preview.
- Can restrict usage to specific Telegram usernames or numeric user IDs.

## What to read next

- [Installation](/getting-started/installation/) for Docker, Compose, and local development.
- [Authentication](/getting-started/authentication/) for Telegram and platform tokens.
- [Configuration](/getting-started/configuration/) for every environment variable.
- [Architecture](/getting-started/architecture/) for the full request flow.
