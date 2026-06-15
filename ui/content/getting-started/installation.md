---
title: "Installation"
description: "Install and run shitpost with Docker, Docker Compose, or Go."
weight: 20
---

## Requirements

- Docker or Podman for the recommended setup.
- Go 1.25+ if you run the bot locally.
- A Telegram bot token from [BotFather](https://t.me/BotFather).
- Credentials for each social platform you want `crosspost` to publish to.

## Option 1: run the image

Use the published image when you want the simplest production path.

```sh
mkdir -p downloads
curl -o .env https://raw.githubusercontent.com/bupd/shitpost/main/.env.example
```

Edit `.env`, then run:

```sh
docker run -d --name shitpost \
  --env-file .env \
  -v ./downloads:/app/downloads \
  registry.goharbor.io/bupd/shitpost:latest
```

Follow logs:

```sh
docker logs -f shitpost
```

Use a versioned tag for production once releases are available. `latest` tracks the newest build from `main`.

## Option 2: Docker Compose

Clone the repo and run the Compose workflow:

```sh
git clone https://github.com/bupd/shitpost.git
cd shitpost
cp .env.example .env
docker compose up --build
```

The service mounts `./downloads` into `/app/downloads` so Telegram media survives container restarts.

## Option 3: Taskfile workflow

If you have [Task](https://taskfile.dev/) installed, the repo exposes common workflows:

```sh
task setup
task up:dry-run
task logs
task down
task doctor
task validate
```

Start with `task up:dry-run`. Your Telegram messages will not be posted; the bot replies with the exact `crosspost` command it would run.

## Option 4: local Go process

Local mode is useful while changing bot behavior.

```sh
cp .env.example .env
task dev:dry-run
```

The local process expects `crosspost` to be available on `PATH` if dry-run is disabled. The Docker image already installs and wraps `crosspost` for you.

## Confirm it works

The bot is ready when logs include:

```text
Authorized as @your_bot_username
```

Send a Telegram message to the bot. In dry-run mode, the reply should start with `DRY RUN: would run`.
