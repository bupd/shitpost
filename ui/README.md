# shitpost docs UI

This folder contains a minimal Hugo documentation site for `shitpost`.

## Local development

```sh
hugo server --source ui
```

## Build

```sh
hugo --source ui --gc --minify
```

## Netlify

Set the Netlify base directory to `ui`. Netlify will read `ui/netlify.toml`, run Hugo, and publish `ui/public`.
