# 🔐 Ephemeral Storage

A lightweight Go service for sharing secrets that self-destruct. Values are **readable only once** and automatically expire after a configurable TTL.

## Features

- **Read-once** — values are deleted immediately after being read
- **TTL expiry** — values auto-expire (default: 60 seconds)
- **Server-generated keys** — base64-encoded `sender.datatype.timestamp`
- **Live dashboard** — real-time storage view powered by SSE
- **API docs** — built-in Swagger UI at `/docs.html`
- **NAIS-ready** — Dockerfile and `nais.yaml` included

## Quick start

```sh
go run ./cmd/server
# → http://localhost:8080
```

## API

### Store a value

```sh
curl -X POST http://localhost:8080/keys \
  -H 'Content-Type: application/json' \
  -d '{"sender":"alice","datatype":"password","value":"s3cret","ttl":60}'
```

Response:

```json
{"key":"YWxpY2UucGFzc3dvcmQuMTcwOTMwMDAwMDAwMA=="}
```

### Read a value (once)

```sh
curl http://localhost:8080/keys/YWxpY2UucGFzc3dvcmQuMTcwOTMwMDAwMDAwMA==
```

Response:

```json
{"value":"s3cret"}
```

A second read returns `404`.

### Live updates (SSE)

```sh
curl http://localhost:8080/events
```

## Docker

```sh
# Build
CGO_ENABLED=0 go build -o server ./cmd/server
docker build -t ephemral-storage .

# Run
docker run -p 8080:8080 ephemral-storage
```

## NAIS deployment

Update `namespace`, `team`, and `ingresses` in `nais.yaml`, then:

```sh
kubectl apply -f nais.yaml
```

## Testing

```sh
go test ./...                                    # all tests
go test ./tests -run TestHTTPInsertAndRead       # single test
```

## Project structure

```
cmd/server/          Entry point and static assets
storage/             Core logic: Store interface, MemoryStore, Service, Handler, SSE hub
tests/               Integration and unit tests
nais.yaml            NAIS application manifest
Dockerfile           Container image definition
```
