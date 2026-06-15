---
title: "Architecture"
description: "How Telegram messages move through shitpost, crosspost, and target platforms."
weight: 50
---

`shitpost` is intentionally thin. It owns the Telegram interface, local media download, environment normalization, and Telegram replies. Publishing is delegated to `crosspost`.

```text
Telegram user
  -> Telegram Bot API
  -> shitpost Go process
  -> crosspost CLI
  -> Bluesky / Mastodon / X / other targets
```

## Runtime pieces

| Piece | Responsibility |
| --- | --- |
| Telegram Bot API | Delivers messages, photos, videos, and documents to the bot. |
| `main.go` | Polls updates, authorizes users, parses messages, downloads files, and starts `crosspost`. |
| `downloads/` | Local working directory for Telegram media files. |
| `crosspost` CLI | Publishes to the selected social platforms. |
| `.env` | Holds Telegram and platform credentials. |
| Docker image | Bundles the Go bot, Bun runtime, and built `crosspost` CLI. |

## Message flow

1. The bot starts and reads `BOT_TOKEN`, `AUTHORIZED_TELEGRAM_USERS`, `CROSSPOST_FLAGS`, and `SHITPOST_DRY_RUN`.
2. The Telegram SDK opens a long-polling update channel.
3. Non-message updates are ignored.
4. If `AUTHORIZED_TELEGRAM_USERS` is set, the sender username or numeric ID must match.
5. Text messages call `crosspost` directly with the message body.
6. Photos are downloaded from Telegram, saved to `downloads/`, and passed to `crosspost` with `--image`.
7. Captions ending with `alt:` are split into clean caption and image alt text.
8. LinkedIn flags are removed from the base target set, then `-l` is added back only when the final post text contains a hashtag.
9. Videos and non-image documents are downloaded, but only caption text is posted by the current CLI path.
10. The bot sends the captured `crosspost` output back to the Telegram chat.

## Docker build flow

The Dockerfile has two stages:

1. A Go builder compiles the static `bot` binary.
2. A Bun runtime clones `crosspost`, installs dependencies, builds it, writes a `/usr/local/bin/crosspost` wrapper, and copies the Go binary.

The final container starts `./bot` from `/app` and stores downloaded media in `/app/downloads`.

## CI and releases

Pull requests and pushes run formatting, vet, build, and tests. Pushes to `main` also build multi-arch images for `linux/amd64` and `linux/arm64`, sign them, publish to GHCR, and copy `latest` to Harbor.

Tags beginning with `v` run GoReleaser, publish release archives, build semver image tags, copy them to Harbor, and sign the images.

## Security model

- The service is self-hosted; messages and downloaded media stay on your machine or server.
- No web UI is exposed by this project.
- The Telegram bot should be treated as a private command surface.
- `AUTHORIZED_TELEGRAM_USERS` should be set for any real account.
- `.env` should never be committed or shared.
