# 2026-04-11 Go Migration Bootstrap

## Context

This session starts a progressive migration path from Rust to Go without breaking the current Rust runtime.

## Decisions

- Keep Rust as production runtime during migration.
- Build Go implementation in a parallel module under `go/`.
- Start from protocol compatibility (external-tagged JSON envelope).
- Add Go tests early to enforce wire format compatibility.

## Implemented

- Added Go module bootstrap:
  - `go/go.mod`
  - `go/cmd/mato/main.go`
  - `go/README.md`
- Added initial protocol compatibility package:
  - `go/protocol/messages.go`
- Added initial Go tests:
  - `go/protocol/messages_test.go`

## Scope of Protocol Port (Initial)

Client variants currently handled in Go:
- Hello
- Spawn
- Input
- Paste
- Resize
- GetScreen
- Subscribe
- GetCwd
- GetIdleStatus

Server variants currently handled in Go:
- Welcome
- Error
- InputModes
- UpdateStatus
- Cwd

## Next Increment

- Add remaining protocol variants (Screen, ScreenDiff, Graphics, IdleStatus, ProcessStatus).
- Add msgpack compatibility tests against Rust fixtures.
- Start Go daemon skeleton that can accept Hello and return Welcome.

## Increment Completed (Daemon Handshake)

- Added Go daemon package with Unix socket support and line-delimited JSON handling:
  - `go/daemon/server.go`
- Added daemon handshake integration test:
  - `go/daemon/server_test.go`
- Extended Go CLI command with modes:
  - `-mode daemon`
  - `-mode hello`
  - optional `-socket` and `-version`
  - file: `go/cmd/mato/main.go`

Validation completed:

- `cd go && go test ./protocol`
- `cd go && go test ./daemon -v`
- `cd go && go test ./cmd/mato`
- Local CLI install and runtime handshake verified:
  - installed to `~/.local/bin/mato-go`
  - daemon started on `/tmp/mato-go-live.sock`
  - hello response: `mato-go hello ok: daemon_version=0.9.6-go socket=/tmp/mato-go-live.sock`
