# Copilot Instructions

This is a Go service for ephemeral key-value storage. Values are readable only once and expire after a TTL.

## Build & Test

```sh
go build ./...
go run ./cmd/server          # start server on :8080
go test ./...
go test ./storage -run TestHTTPInsertAndRead  # single test
```

## API

- `POST /keys` — `{"key":"k","value":"v","ttl":60}` → 201
- `GET /keys/{key}` — returns `{"value":"v"}` → 200 (first read), 404 (after)

## Architecture

- `storage.Store` — interface for storage backends (`Set`, `GetAndDelete`)
- `storage.MemoryStore` — in-memory implementation (thread-safe)
- `storage.Service` — business logic layer exposing `InsertKeyValue` and `ReadValue`
- `storage.Handler` — HTTP handlers wrapping the Service
- `cmd/server` — main entry point, wires everything together

New backends should implement the `Store` interface.

## Conventions

- Language: Go
- Storage backends go in the `storage` package and implement the `Store` interface
- Values are read-once: `GetAndDelete` must atomically retrieve and remove the entry
- All store implementations must be safe for concurrent use
