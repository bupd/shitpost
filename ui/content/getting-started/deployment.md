---
title: "Deployment"
description: "Run the bot on a server and deploy these docs on Netlify."
weight: 70
---

## Bot deployment

For a server, prefer the published container image and a persistent `downloads` directory.

```sh
mkdir -p /opt/shitpost/downloads
cd /opt/shitpost
curl -o .env https://raw.githubusercontent.com/bupd/shitpost/main/.env.example
```

Edit `/opt/shitpost/.env`, then run:

```sh
docker run -d --name shitpost \
  --restart unless-stopped \
  --env-file /opt/shitpost/.env \
  -v /opt/shitpost/downloads:/app/downloads \
  registry.goharbor.io/bupd/shitpost:latest
```

Check logs:

```sh
docker logs -f shitpost
```

## Docker Compose deployment

The repo includes `docker-compose.yml`:

```sh
cp .env.example .env
docker compose up --build -d
docker compose logs -f shitpost-bot
```

The Compose service is named `shitpost-bot` and the container is named `shitpost-engine`.

## Image registries

Published images are available from:

```sh
docker pull registry.goharbor.io/bupd/shitpost:latest
docker pull ghcr.io/bupd/shitpost:latest
```

For production, use a versioned tag when one is available.

## Docs deployment on Netlify

This `ui` folder is a Hugo site. In Netlify:

1. Connect the Git repository.
2. Set the base directory to `ui`.
3. Keep the build command as `hugo --gc --minify`.
4. Keep the publish directory as `public`.

`ui/netlify.toml` already contains those settings.

## Local docs preview

```sh
hugo server --source ui
```

Build the static site:

```sh
hugo --source ui --gc --minify
```
