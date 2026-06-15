# shitpost

<p align="center">
  <img src="ui/static/icon.svg" alt="shitpost" width="96" height="96">
</p>

<p align="center">
  <strong>Post once from Telegram. Publish everywhere.</strong>
</p>

<p align="center">
  <a href="https://github.com/bupd/shitpost/actions/workflows/ci.yml"><img src="https://github.com/bupd/shitpost/actions/workflows/ci.yml/badge.svg?branch=main" alt="CI"></a>
  <a href="https://shitpost.bupd.xyz"><img src="https://img.shields.io/badge/docs-shitpost.bupd.xyz-2874e4" alt="Documentation"></a>
  <a href="https://github.com/bupd/shitpost"><img src="https://img.shields.io/badge/self--hosted-yes-24292f" alt="Self-hosted"></a>
</p>

`shitpost` is a lightweight, self-hosted Telegram bot for crossposting to your social accounts. Send a message, photo, caption, or alt text to a private Telegram bot and let `crosspost` publish it to the platforms you configured.

It is built for people who want the speed of posting from chat without handing their tokens, drafts, or media to a hosted scheduler.

## Why use it?

- One private Telegram chat becomes your posting interface.
- No dashboard, queue, browser tab, or SaaS account required.
- Text and image posts work from mobile or desktop Telegram.
- Alt text is supported with a simple `alt:` caption suffix.
- Dry-run mode previews the exact `crosspost` command before publishing.
- Docker images are available for `linux/amd64` and `linux/arm64`.

## How it works

```text
Telegram -> shitpost bot -> crosspost CLI -> Bluesky / Mastodon / X / more
```

`shitpost` owns the Telegram side: authorization, message parsing, media downloads, environment normalization, and replies with logs. Publishing is delegated to the `crosspost` CLI.

## Quick start

Create an env file:

```sh
curl -o .env https://raw.githubusercontent.com/bupd/shitpost/main/.env.example
```

Edit `.env`, then start the bot:

```sh
mkdir -p downloads
docker run -d --name shitpost \
  --env-file .env \
  -v ./downloads:/app/downloads \
  registry.goharbor.io/bupd/shitpost:latest
```

Watch the logs:

```sh
docker logs -f shitpost
```

The bot is ready when logs include:

```text
Authorized as @your_bot_username
```

## Recommended first run

Start in dry-run mode before posting for real:

```dotenv
SHITPOST_DRY_RUN=1
AUTHORIZED_TELEGRAM_USERS=your_telegram_username
CROSSPOST_FLAGS=-bmt
```

Send a Telegram message to the bot. It will reply with the command it would run instead of publishing.

## Documentation

Full setup docs live at [shitpost.bupd.xyz](https://shitpost.bupd.xyz).

Start there for:

- Installation with Docker, Compose, or local Go.
- Telegram bot creation and private access control.
- Bluesky, Mastodon, and X credential setup.
- Every supported environment variable and alias.
- Architecture, deployment, usage, and troubleshooting.

## Environment variables

| Variable | Required | Description |
| --- | --- | --- |
| `BOT_TOKEN` | Yes | Telegram bot token from BotFather. |
| `AUTHORIZED_TELEGRAM_USERS` | Recommended | Comma-separated Telegram usernames or numeric user IDs allowed to use the bot. |
| `CROSSPOST_FLAGS` | No | Flags passed to `crosspost`; defaults to `-bmt`. |
| `SHITPOST_DRY_RUN` | No | Set to `1`, `true`, or `yes` to preview commands without posting. |
| `AUTH_TOKEN` | No | X `auth_token` cookie value used by the emusks-backed `crosspost` strategy. |
| `TWITTER_AUTH_TOKEN` | No | Explicit X auth token; overrides `AUTH_TOKEN`. |
| `BLUESKY_HOST` | No | Bluesky host, usually `bsky.social`. |
| `BLUESKY_IDENTIFIER` | No | Bluesky handle or email. |
| `BLUESKY_PASSWORD` | No | Bluesky app password. |
| `MASTODON_HOST` | No | Mastodon instance URL. |
| `MASTODON_ACCESS_TOKEN` | No | Mastodon access token. |

Run `task doctor` to check which secrets are present without printing their values.

## Taskfile workflow

```sh
task setup              # create .env and download Go deps
task up                 # run in Docker Compose
task up:dry-run         # run without posting
task up:detached        # run in the background
task logs               # follow container logs
task doctor             # validate .env shape without leaking secrets
task validate           # gofmt, go vet, go test, go build
```

## Media and alt text

Text messages are posted as text. Photos are downloaded, attached, and posted with their caption.

Add alt text by ending a caption with a final `alt:` line:

```text
new deploy view from the homelab
alt: A terminal window showing a successful Docker deployment
```

Videos and non-image documents are downloaded, but the current `crosspost` CLI path posts caption text only.

## Images

```sh
docker pull registry.goharbor.io/bupd/shitpost:latest
docker pull ghcr.io/bupd/shitpost:latest
```

Use versioned tags for production when available. `latest` tracks the newest build from `main`.

## Security

- Set `AUTHORIZED_TELEGRAM_USERS` for any real deployment.
- Keep `.env` out of Git and backups you do not control.
- Treat Telegram bot tokens and platform tokens like passwords.
- Mount `downloads/` somewhere persistent if you want retained media.

## Related projects

- [crosspost](https://github.com/humanwhocodes/crosspost) powers the platform publishing.
- [Postiz](https://github.com/gitroomhq/postiz-app) is a full-featured social media scheduler.
- [Buffer](https://buffer.com) is a commercial social media management product.

## License

See [LICENSE](LICENSE).
