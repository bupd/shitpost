---
title: "Configuration"
description: "Environment variables, target flags, aliases, and safe validation."
weight: 40
---

## Minimal `.env`

Start with the example file:

```sh
cp .env.example .env
```

The smallest useful dry-run config is:

```dotenv
BOT_TOKEN=123456789:telegram-secret
AUTHORIZED_TELEGRAM_USERS=your_username
CROSSPOST_FLAGS=-bmt
SHITPOST_DRY_RUN=1
```

`CROSSPOST_FLAGS=-bmt` tells `crosspost` to post to Bluesky, Mastodon, and Twitter/X. Change it to match the destinations you actually configured.

## Core variables

| Variable | Required | Purpose |
| --- | --- | --- |
| `BOT_TOKEN` | Yes | Telegram bot token from BotFather. |
| `AUTHORIZED_TELEGRAM_USERS` | Recommended | Comma-separated usernames or numeric IDs allowed to use the bot. Empty means public. |
| `CROSSPOST_FLAGS` | No | Flags passed directly to `crosspost`; defaults to `-bmt`. |
| `SHITPOST_DRY_RUN` | No | Set to `1`, `true`, or `yes` to preview commands without posting. |
| `CROSSPOST_REPO` | No | Source repo used when building the Docker image. |
| `CROSSPOST_REF` | No | Branch, tag, or ref used when building the Docker image. |

## Platform variables

| Variable | Platform | Purpose |
| --- | --- | --- |
| `BLUESKY_HOST` | Bluesky | Usually `bsky.social`. |
| `BLUESKY_IDENTIFIER` | Bluesky | Handle or email. |
| `BLUESKY_PASSWORD` | Bluesky | App password. |
| `MASTODON_HOST` | Mastodon | Instance URL, such as `https://mastodon.social`. |
| `MASTODON_ACCESS_TOKEN` | Mastodon | Access token with posting permission. |
| `MASTODON_CLIENT_KEY` | Mastodon | Optional client key for crosspost flows that need it. |
| `MASTODON_CLIENT_SECRET` | Mastodon | Optional client secret for crosspost flows that need it. |
| `AUTH_TOKEN` | X | Alias for the X `auth_token` cookie. |
| `TWITTER_AUTH_TOKEN` | X | Explicit X `auth_token`; overrides `AUTH_TOKEN`. |
| `TWITTER_AUTH_CLIENT` | X | Optional client identity understood by the emusks-backed strategy. |
| `TWITTER_GRAPHQL_ENDPOINT` | X | Optional endpoint profile understood by the emusks-backed strategy. |
| `TWITTER_PROXY` | X | Optional proxy for X requests. |
| `TWITTER_API_CONSUMER_KEY` | X | Official API consumer key fallback. |
| `TWITTER_API_CONSUMER_SECRET` | X | Official API consumer secret fallback. |
| `TWITTER_ACCESS_TOKEN_KEY` | X | Official API access token fallback. |
| `TWITTER_ACCESS_TOKEN_SECRET` | X | Official API access token secret fallback. |

## Alias normalization

Before `shitpost` starts `crosspost`, it copies legacy aliases into the variable names `crosspost` expects when the target variable is empty.

| Target | Accepted aliases |
| --- | --- |
| `TWITTER_AUTH_TOKEN` | `AUTH_TOKEN` |
| `TWITTER_API_CONSUMER_KEY` | `consumer_key`, `TWITTER_CONSUMER_KEY` |
| `TWITTER_API_CONSUMER_SECRET` | `consumer_key_secret`, `TWITTER_CONSUMER_SECRET` |
| `TWITTER_ACCESS_TOKEN_KEY` | `access_token`, `access_token_key`, `TWITTER_ACCESS_TOKEN` |
| `TWITTER_ACCESS_TOKEN_SECRET` | `access_token_secret`, `TWITTER_ACCESS_SECRET` |

## Validate without leaking secrets

Run the doctor task when posts fail or before deploying a new `.env`:

```sh
task doctor
```

The doctor script prints whether each key is present and the value length. It does not print secret values.

## Recommended first run

```dotenv
SHITPOST_DRY_RUN=1
CROSSPOST_FLAGS=-bmt
```

Send one text message and one image with a caption. If the Telegram replies show the expected command, set `SHITPOST_DRY_RUN=0` and restart the service.
