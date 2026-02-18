# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CLI tool for downloading Zoom meeting recordings via the Zoom API. Built in Go using Kong for CLI parsing and zerolog for structured logging. Currently in early/template stage — the `Options.Run()` method is a stub.

## Build & Development Commands

- **Build:** `make` or `make compile` (binary output: `./zoom-recordings`)
- **Test:** `make test` (runs `go test -v ./...`)
- **Run single test:** `go test -v -run TestName ./package/...`
- **Vet:** `make vet` (runs `go vet ./...`)
- **Install:** `make install`
- **Install lint tools:** `make tools` (staticcheck, gocritic, gosec)

Version is injected at build time via `-ldflags` from `git describe --tags`.

## Architecture

- **`main.go`** — Entry point. Sets up Kong CLI parsing, binds `context.Context`, `zerolog.Logger`, and `Options` into Kong's dependency injection, then calls `kongContext.Run()`.
- **`program/`** — Core package:
  - `program.go` — `Options` struct (Kong-annotated CLI flags), `Parse()`, `Run()`, logging setup. Kong's `AfterApply` hook initializes logging before command execution.
  - `version.go` — `Version` variable set by linker flags at build time.
  - `program_test.go` — Tests using `testify/assert`, `monkey` (patching `os.Exit`), and `go-capturer`.

# API References
## Zoom API Reference

- **Docs:** https://developers.zoom.us/docs/api/
- **Context7 library:** `/zoom/server-to-server-oauth-starter-api` (use with Context7 MCP plugin for recordings, OAuth, and user settings endpoints)

## Kong CLI Pattern

- **Docs:** https://github.com/alecthomas/kong
- **Examples:** https://github.com/alecthomas/kong/tree/master/_examples

Commands are added as sub-structs on `Options` with Kong `cmd:""` tags. Kong bindings (`kongContext.BindTo`, `kongContext.Bind`) make `context.Context`, `zerolog.Logger`, and `Options` available as injectable arguments to any command's `Run()` method.
