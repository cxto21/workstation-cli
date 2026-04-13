# Mato Go (Bootstrap)

This directory contains the progressive Go migration for Mato.

Current scope:
- Go module and command bootstrap
- IPC protocol compatibility layer (initial subset)
- Compatibility tests for external-tagged JSON envelope used by Rust
- Minimal Unix socket daemon with Hello -> Welcome handshake

Not yet included:
- Daemon runtime
- TUI client runtime
- PTY provider/emulator parity

## Run tests

```bash
cd go
go test ./...
```

## Run daemon and hello

```bash
cd go
go run ./cmd/mato -mode daemon
```

In another terminal:

```bash
cd go
go run ./cmd/mato -mode hello
```
