# Deploy Scaffold

## Local

- `local/docker-compose.yml`: minimal Redis container for local development

## Production

- `prod/Caddyfile`: reverse proxy starter for the API domain
- `prod/.env.example`: sample environment variables for the Go API

This folder is intentionally small for now. It gives us just enough structure to wire local development and a low-cost VPS deployment without locking us into a heavier platform too early.
