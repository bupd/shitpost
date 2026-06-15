---
title: "Usage"
description: "Send text, images, captions, and alt text from Telegram."
weight: 60
---

## Text posts

Send a normal text message to your Telegram bot. `shitpost` passes that text as the final argument to `crosspost`.

```text
shipping a small self-hosted crossposter today
```

In dry-run mode, the bot replies with a command preview. In live mode, it replies with `crosspost` stdout and stderr.

## Image posts

Send a photo with an optional caption. Telegram provides multiple photo sizes; `shitpost` picks the largest one, saves it under `downloads/`, and passes it with `--image`.

```text
caption: a tiny bot doing useful work
```

## Alt text

Add alt text by ending the caption with a final line that starts with `alt:`.

```text
new deploy view from the homelab
alt: A terminal window showing a successful Docker deployment
```

The posted caption becomes `new deploy view from the homelab`. The alt text is passed separately to `crosspost`.

## Videos and documents

Telegram videos and documents are downloaded so the bot can acknowledge and inspect them. The current installed `crosspost` CLI path only attaches images, so video and non-image document messages post caption text only.

If there is no caption, the bot replies with a warning instead of posting an empty update.

## Dry-run mode

Dry-run mode is the safest way to test credentials, target flags, captions, and alt text.

```dotenv
SHITPOST_DRY_RUN=1
```

Restart the service after changing `.env`. The startup logs should include:

```text
Dry-run mode enabled. Messages will not be posted.
```

## Target selection

`CROSSPOST_FLAGS` is passed directly to `crosspost`.

```dotenv
CROSSPOST_FLAGS=-bmt
```

Use the flags supported by the `crosspost` version you build into the image. The default project configuration targets Bluesky, Mastodon, and Twitter/X.

## LinkedIn hashtag gate

LinkedIn is added only when the outgoing post contains a hashtag.

```text
short note for the timeline
```

This posts to the base targets from `CROSSPOST_FLAGS`, usually Bluesky, Mastodon, and X.

```text
shipping the new docs site #golang
```

This posts to the base targets and LinkedIn. The rule is enforced even if `CROSSPOST_FLAGS` accidentally includes `-l`, `--linkedin`, or `-bmtl`.
