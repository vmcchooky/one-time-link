# one-time-link

`one-time-link` is a one-time secret sharing application designed for portfolio use first, with a path toward production hardening later.

## Current focus

The repository is currently set up for Milestone 1:

- product and deployment documentation is in place
- the frontend has a React + TypeScript shell
- the backend has a Go API shell
- local Redis is wired through Docker Compose
- deployment docs target `quorix.io.vn` with a low-cost VPS primary and Oracle Cloud standby

## Repository structure

```text
backend/
  cmd/api/
  internal/
frontend/
  web-app/
deploy/
  local/
  prod/
docs/
scripts/
```

## Local development

### 1. Start Redis

```bash
docker compose -f deploy/local/docker-compose.yml up -d
```

### 2. Run the Go API

```bash
go run ./backend/cmd/api
```

The API listens on `http://localhost:8080` by default.

### 3. Run the frontend

```bash
cd frontend/web-app
npm install
npm run dev
```

The frontend runs on `http://localhost:5173` by default.

## Key docs

- [docs/README.md](/D:/Go/duan/one-time-link/docs/README.md)
- [docs/product-spec/one-time-link-requirements.md](/D:/Go/duan/one-time-link/docs/product-spec/one-time-link-requirements.md)
- [docs/product-spec/one-time-link-architecture.md](/D:/Go/duan/one-time-link/docs/product-spec/one-time-link-architecture.md)
- [docs/product-spec/one-time-link-milestones.md](/D:/Go/duan/one-time-link/docs/product-spec/one-time-link-milestones.md)
- [docs/deployment/deployment-decision-summary.md](/D:/Go/duan/one-time-link/docs/deployment/deployment-decision-summary.md)

## Milestone 1 done in code

- monorepo layout
- frontend shell
- backend shell
- local Redis compose file
- shared HTTP contract document
- health endpoint and structured request logging

## Next step

Milestone 2 is to implement the create-secret flow with client-side encryption and a real `POST /api/secrets` backend path.
