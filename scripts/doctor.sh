#!/usr/bin/env sh
set -eu

env_file="${1:-.env}"

if [ ! -f "$env_file" ]; then
  printf 'missing %s\n' "$env_file"
  exit 1
fi

set -a
# shellcheck disable=SC1090
case "$env_file" in
  */*) . "$env_file" ;;
  *) . "./$env_file" ;;
esac
set +a

check_key() {
  key="$1"
  value=$(printenv "$key" 2>/dev/null || true)
  if [ -n "$value" ]; then
    printf '%s=present length=%s\n' "$key" "${#value}"
  else
    printf '%s=missing\n' "$key"
  fi
}

check_alias() {
  target="$1"
  shift
  if [ -n "$(printenv "$target" 2>/dev/null || true)" ]; then
    return
  fi

  for alias in "$@"; do
    value=$(printenv "$alias" 2>/dev/null || true)
    if [ -n "$value" ]; then
      printf '%s can be derived from %s\n' "$target" "$alias"
      return
    fi
  done
}

printf 'Telegram\n'
check_key BOT_TOKEN
check_key AUTHORIZED_TELEGRAM_USERS
check_key CROSSPOST_FLAGS
check_key SHITPOST_DRY_RUN

printf '\nBluesky\n'
check_key BLUESKY_HOST
check_key BLUESKY_IDENTIFIER
check_key BLUESKY_PASSWORD

printf '\nMastodon\n'
check_key MASTODON_HOST
check_key MASTODON_ACCESS_TOKEN

printf '\nTwitter/X\n'
check_key AUTH_TOKEN
check_key TWITTER_AUTH_TOKEN
check_key TWITTER_API_CONSUMER_KEY
check_key TWITTER_API_CONSUMER_SECRET
check_key TWITTER_ACCESS_TOKEN_KEY
check_key TWITTER_ACCESS_TOKEN_SECRET
check_alias TWITTER_AUTH_TOKEN AUTH_TOKEN
check_alias TWITTER_API_CONSUMER_KEY consumer_key TWITTER_CONSUMER_KEY
check_alias TWITTER_API_CONSUMER_SECRET consumer_key_secret TWITTER_CONSUMER_SECRET
check_alias TWITTER_ACCESS_TOKEN_KEY access_token access_token_key TWITTER_ACCESS_TOKEN
check_alias TWITTER_ACCESS_TOKEN_SECRET access_token_secret TWITTER_ACCESS_SECRET

access_token=$(printenv TWITTER_ACCESS_TOKEN_KEY 2>/dev/null || true)
if [ -n "$access_token" ] && ! printf '%s' "$access_token" | grep -q '-'; then
  printf '\nwarning: TWITTER_ACCESS_TOKEN_KEY does not look like a typical OAuth 1.0a access token. If X returns 401, confirm this is Access Token, not OAuth2 Client ID.\n'
fi
