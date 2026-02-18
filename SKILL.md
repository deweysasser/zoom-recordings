---
name: zoom-recordings
description: Use when the user asks to download, list, or manage Zoom meeting recordings, or when automating Zoom recording retrieval as part of a larger workflow.
---

# zoom-recordings CLI

## Overview

`zoom-recordings` is a Go CLI tool in this repository for downloading Zoom meeting recordings via OAuth. Use it to authenticate with Zoom, list available recordings, and download them to local files.

## Prerequisites

The user needs a Zoom OAuth App with a **Client ID** and **Client Secret**. The redirect URL must be `http://localhost:8085/oauth/callback` (or match the configured `--callback-port`).

## Authentication

Authentication must happen first. It opens a browser for OAuth consent and saves a token to `~/.zoom-recordings/token.json`.

```sh
zoom-recordings login --client-id=ID --client-secret=SECRET
```

Or with environment variables already set (`ZOOM_CLIENT_ID`, `ZOOM_CLIENT_SECRET`):

```sh
zoom-recordings login
```

The token is reused automatically by `list` and `download`. Re-run `login` if commands fail with auth errors.

## Listing Recordings

```sh
# Last 24 hours (default)
zoom-recordings list

# Specific date range
zoom-recordings list --from 2026-01-01 --to 2026-02-01
```

## Downloading Recordings

```sh
# Download to current directory (default)
zoom-recordings download

# Custom output directory and date range
zoom-recordings download --output-dir ./my-recordings --from 2026-01-01 --to 2026-02-01
```

Files are named `{date}_{topic}_{recording-type}.{ext}`. Downloads are idempotent — existing files with matching sizes are skipped.

## Common Workflows

**First-time setup:** Run `login`, then `download`. Suggest the user set `ZOOM_CLIENT_ID` and `ZOOM_CLIENT_SECRET` as environment variables to avoid repeating them.

**Routine backup:** Run `download` with no flags to grab the last 24 hours. Safe to run repeatedly.

**Specific meeting range:** Use `--from` and `--to` on either `list` (to preview) or `download`.

## When Things Go Wrong

| Symptom | Fix |
|---------|-----|
| "not authenticated" error | Run `zoom-recordings login` first |
| Auth errors after token was saved | Token may be expired; re-run `login` |
| No recordings found | Check date range; default is last 24 hours only |
| Download hangs | Check network; Zoom download URLs require internet access |
