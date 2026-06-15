---
title: "Authentication"
description: "Create the Telegram bot and collect the platform credentials that crosspost needs."
weight: 30
---

## Telegram bot token

1. Open [@BotFather](https://t.me/BotFather) in Telegram.
2. Send `/newbot`.
3. Pick a display name and bot username.
4. Copy the token BotFather returns.
5. Put it in `.env` as `BOT_TOKEN`.

```dotenv
BOT_TOKEN=123456789:telegram-secret
```

Treat the bot token like a password. Anyone with it can control the bot API for that bot.

## Private bot access

Set `AUTHORIZED_TELEGRAM_USERS` so only trusted Telegram accounts can trigger posts.

```dotenv
AUTHORIZED_TELEGRAM_USERS=alice,123456789
```

Accepted values are Telegram usernames without `@`, usernames with `@`, or numeric Telegram user IDs. Separate multiple users with commas.

Leaving this variable empty means anyone who can message the bot can ask it to post. That is only safe for a throwaway test bot.

## Bluesky credentials

Use an app password instead of your account password.

1. Open Bluesky settings.
2. Go to app passwords.
3. Create a new app password for `shitpost`.
4. Set your handle or email as `BLUESKY_IDENTIFIER`.
5. Set the app password as `BLUESKY_PASSWORD`.

```dotenv
BLUESKY_HOST=bsky.social
BLUESKY_IDENTIFIER=you.bsky.social
BLUESKY_PASSWORD=xxxx-xxxx-xxxx-xxxx
```

## Mastodon credentials

Create an application from your Mastodon instance preferences.

1. Open your Mastodon instance in a browser.
2. Go to Preferences, then Development.
3. Create a new application with write scopes.
4. Copy the access token.
5. Set the instance URL and token in `.env`.

```dotenv
MASTODON_HOST=https://mastodon.social
MASTODON_ACCESS_TOKEN=secret-token
```

Some `crosspost` flows may also use `MASTODON_CLIENT_KEY` and `MASTODON_CLIENT_SECRET`; keep them available if your target command needs them.

## LinkedIn credentials

LinkedIn posts are intentionally gated by hashtags. Configure the credential once, then `shitpost` will add LinkedIn only for posts whose final text or caption contains a hashtag.

1. Go to [LinkedIn Developers](https://www.linkedin.com/developers/).
2. Create an app and verify the LinkedIn page or member identity you want to post from.
3. Request access to "Share on LinkedIn" and "Sign in with LinkedIn using OpenID Connect".
4. Open [OAuth 2.0 Tools](https://www.linkedin.com/developers/tools/oauth) and create a token for your app.
5. Include the `openid`, `profile`, and `w_member_social` scopes.
6. Set the access token in `.env`.

```dotenv
LINKEDIN_ACCESS_TOKEN=linkedin-access-token
```

LinkedIn tokens can expire. If hashtagged posts start failing while other networks work, create a fresh token and restart the bot.

## Twitter / X credentials

`shitpost` supports the `crosspost` Twitter strategies and normalizes several legacy aliases before invoking `crosspost`.

For the emusks-backed strategy, provide the X `auth_token` cookie value:

```dotenv
AUTH_TOKEN=your-x-auth-token-cookie
```

You can also set the explicit variable:

```dotenv
TWITTER_AUTH_TOKEN=your-x-auth-token-cookie
```

If both are set, `TWITTER_AUTH_TOKEN` wins.

For official API fallback paths, use OAuth 1.0a credentials:

```dotenv
TWITTER_API_CONSUMER_KEY=
TWITTER_API_CONSUMER_SECRET=
TWITTER_ACCESS_TOKEN_KEY=
TWITTER_ACCESS_TOKEN_SECRET=
```

To find the X `auth_token` cookie, sign in to X in a browser you control, open developer tools, inspect cookies for `x.com`, and copy the value named `auth_token`. Do not paste the full cookie header; only the token value belongs in `.env`.
