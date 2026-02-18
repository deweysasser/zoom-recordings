# zoom-recordings

A CLI tool for downloading Zoom meeting recordings via the Zoom API.

## Prerequisites

You need a Zoom OAuth App to authenticate. Create one at the [Zoom App Marketplace](https://marketplace.zoom.us/):

1. Create a **User-managed OAuth** app
2. Set the redirect URL to `http://localhost:8085/oauth/callback`
3. Add the `recording:read` scope
4. Note your **Client ID** and **Client Secret**

## Installation

```sh
make install
```

Or build locally:

```sh
make
# binary: ./zoom-recordings
```

## Configuration

Pass your OAuth credentials via flags or environment variables:

```sh
export ZOOM_CLIENT_ID=your_client_id
export ZOOM_CLIENT_SECRET=your_client_secret
```

The callback port defaults to `8085` and can be changed with `--callback-port` or `ZOOM_CALLBACK_PORT`.

## Usage

### Authenticate

Log in via your browser. This saves a token to `~/.zoom-recordings/token.json` that is reused for subsequent commands.

```sh
zoom-recordings login
```

### List recordings

```sh
# List recordings from the last 24 hours (default)
zoom-recordings list

# List recordings in a specific date range
zoom-recordings list --from 2026-01-01 --to 2026-02-01
```

### Download recordings

```sh
# Download to ./recordings (default)
zoom-recordings download

# Download to a specific directory with a date range
zoom-recordings download --output-dir ./my-recordings --from 2026-01-01 --to 2026-02-01
```

Downloaded files are named `{date}_{topic}_{recording-type}.{ext}`, e.g. `2026-02-17_team-standup_shared_screen_with_speaker_view.mp4`.

Downloads are idempotent — files that already exist with the correct size are skipped.

## Claude Code Integration

This project includes a [Claude Code skill](SKILL.md) that teaches Claude how to use the `zoom-recordings` CLI. When working in this repository with Claude Code, it can guide you through authentication, listing, and downloading recordings.

## Building

```sh
make              # build binary
make test         # run tests
make vet          # run go vet
make install      # install to $GOPATH/bin
```
