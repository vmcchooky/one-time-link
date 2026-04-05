# Backend Scaffold

This backend starts with one Go API process so the team can move fast cheaply.

## Current scope

- `cmd/api`: single deployable API entrypoint
- `internal/config`: environment-based config loading
- `internal/httpapi`: HTTP routing and handler layer
- `internal/secret`: placeholder service boundary for secret lifecycle logic

## Why this shape

The code is scaffolded to keep clean domain boundaries now, while still letting deployment stay simple on one VPS.

Later, the same boundaries can be split into dedicated services without rewriting the product model.

## Run locally

```bash
go run ./backend/cmd/api
```
