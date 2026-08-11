# Merkadapp API

REST API for a personal grocery management and expense tracking app. Built as a hands-on project to learn Go, following a layered architecture.

**Live:** [merkadapp.onrender.com](https://merkadapp.onrender.com)

Pairs with [merkadapp_frontend](https://github.com/raulito1500/merkadapp_frontend) — the React SPA that consumes this API.

## Stack

- **Go 1.22**
- **Gin** — HTTP router and middleware
- **MongoDB** — via the official Go driver
- **Gorilla WebSocket** — real-time notifications channel
- **gin-contrib/cors** — CORS configuration

## Architecture

The codebase follows a strict handlers → services → repositories separation, and each domain (`bill`, `market_list`, `product`, `notification`) owns its full stack independently — no shared service layer between domains. All routes are registered in one place (`server/api.go`), which wires each domain's repository → service → handler and mounts its routes on the Gin engine.

There is currently **no authentication/authorization layer** on this API — every route is open. (A `JWT_SECRET` config value exists but is unused dead code, unrelated to the frontend's Firebase login, which this API does not enforce.)

## API Reference

There's no live interactive documentation for this API yet. The routes and their exact paths are registered in [`server/api.go`](server/api.go); request/response shapes for each domain live in its `handlers/` and `entities/` packages (see [Where to find things](#where-to-find-things) below).

## Where to find things

| Looking for... | Go to |
|---|---|
| Route registration (every endpoint, wired to its handler) | `server/api.go` |
| Bills — HTTP layer, business logic (totals, merging, recommendations), data access | `src/bill/{handlers,services,repository,models,entities}` |
| Market lists — same layers | `src/market_list/{handlers,services,repository,models,entities}` |
| Products — same layers | `src/product/{handlers,services,repository,models}` |
| WebSocket upgrade / live push notifications | `src/notification/handlers` |
| Config loading (env vars) | `config/` |
| MongoDB connection setup | `database/` |

## Getting Started

**Prerequisites:** Go 1.22+ and a MongoDB instance (local or [Atlas free tier](https://www.mongodb.com/cloud/atlas))

```bash
git clone https://github.com/raulito1500/merkadapp.git
cd merkadapp

cp .env.example .env
# edit .env and fill in your MongoDB credentials — .env.example documents every variable

go mod download
go run .
```

The server starts at `http://localhost:8080`.

## License

[PolyForm Noncommercial 1.0.0](LICENSE) — free to use, copy, modify and distribute for noncommercial purposes (personal, educational, portfolio). Commercial use requires permission from the author.
