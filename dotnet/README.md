# 🔐 Ephemeral Storage

An ASP.NET Core service for sharing secrets that self-destruct. Values are **readable only once** and automatically expire after a configurable TTL.

## Features

- **Read-once** — values are deleted immediately after being read
- **TTL expiry** — values auto-expire (default: 60 seconds)
- **Server-generated keys** — base64-encoded `sender.datatype.timestamp`
- **Live dashboard** — real-time storage view powered by SSE
- **API docs** — built-in Swagger UI at `/docs.html`

## Quick start

```sh
dotnet run --project EphemralStorage
# → http://localhost:5000
```

## API

### Store a value

```sh
curl -X POST http://localhost:5000/keys \
  -H 'Content-Type: application/json' \
  -d '{"sender":"alice","datatype":"password","value":"s3cret","ttl":60}'
```

Response:

```json
{"key":"YWxpY2UucGFzc3dvcmQuMTcwOTMwMDAwMDAwMA=="}
```

### Read a value (once)

```sh
curl http://localhost:5000/keys/YWxpY2UucGFzc3dvcmQuMTcwOTMwMDAwMDAwMA==
```

Response:

```json
{"value":"s3cret"}
```

A second read returns `404`.

### Live updates (SSE)

```sh
curl http://localhost:5000/events
```

## Testing

```sh
dotnet test                                          # all tests
dotnet test --filter TestInsertAndRead               # single test
dotnet test --filter "PerformanceTests" -c Release   # benchmarks
```

## Project structure

```
EphemralStorage/          Web API entry point and static assets
EphemralStorage.Core/     Core abstractions: IStore interface, EphemeralService
EphemralStorage.Memory/   In-memory IStore implementation
EphemralStorage.Tests/    Unit, integration, and performance tests
```
