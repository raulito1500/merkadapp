# Merkadapp API

REST API for grocery planning and shared expense tracking. It imports Colombian electronic invoices, turns receipt history into price history, and uses that to tell you what a shopping list will probably cost before you go.

**Live:** [merkadapp.onrender.com](https://merkadapp.onrender.com)

Pairs with [merkadapp_frontend](https://github.com/raulito1500/merkadapp_frontend) — the React SPA that consumes this API.

## What's interesting in here

Most of this is ordinary CRUD. Three parts aren't:

**Electronic invoice import** (`src/bill/services/billService.go`). Colombian DIAN invoices arrive as a UBL `AttachedDocument` that wraps the real invoice XML inside a CDATA block — so importing one means unwrapping a document, extracting and re-parsing the payload, then normalizing free-text line descriptions into something comparable: pulling quantity and unit out of strings like `LECHE ENTERA 1000ML`, converting grams to kilos, and cleaning the trailing markers supermarkets add. Upload the XML, get a structured receipt.

**Price history by correlated aggregation** (`src/market_list/repository/marketListMongoRepository.go`). For every item on a shopping list, the list query resolves what that product last cost and where it was last bought, by correlating each item against all prior receipts inside a single `$lookup` sub-pipeline, then summing the result into an estimated total for the list. It's the query the whole app is built around.

**Catalog-gap recommendations** (`RecommendedProducts`). Line items that never got linked to a catalogued product, but show up on three or more receipts, surface as suggestions — the things you clearly keep buying that nobody ever bothered to add as products.

## Stack

- **Go 1.22**
- **Gin** — HTTP router and middleware
- **MongoDB** — via the official Go driver, aggregation pipelines written by hand
- **Gorilla WebSocket** — connection upgrade for a notifications channel (see *Status* below)

## Architecture

Strict handlers → services → repositories separation, with each domain (`bill`, `market_list`, `product`, `notification`) owning its full stack independently — no shared service layer between domains. All routes are registered in one place, [`server/api.go`](server/api.go), which wires each domain's repository → service → handler and mounts its routes on the Gin engine.

## Status and known gaps

This started as a personal project and still runs as one. It works and it's deployed, but it has never had to survive anyone but me, and the gaps show:

- **No authentication or authorization.** Every route is open. A `JWT_SECRET` config value exists but is dead code, unrelated to the frontend's Firebase login, which this API does not enforce. This is the first thing that would have to change before it was worth anything to anyone else.
- **CORS allows all origins**, which is a deliberate convenience that pairs badly with the point above.
- **The WebSocket endpoint is a stub.** It upgrades the connection and then writes a fixed message on a timer. The channel exists; nothing meaningful flows through it yet.
- **Test coverage is one file.** Validators only.
- **MongoDB was arguably the wrong choice.** Receipts, lists and products are relational — a product appears on lists, gets bought on receipts, accumulates a price history. A fair amount of the aggregation work above would be shorter and more consistent as SQL against Postgres with real foreign keys.

## API reference

There's no generated documentation yet. Routes and their exact paths are registered in [`server/api.go`](server/api.go); request and response shapes for each domain live in its `handlers/` and `entities/` packages.

## Where to find things

| Looking for... | Go to |
|---|---|
| Route registration (every endpoint, wired to its handler) | `server/api.go` |
| Bills — HTTP layer, totals, merging, invoice import, recommendations | `src/bill/{handlers,services,repository,models,entities}` |
| Market lists — same layers, plus the price-history aggregation | `src/market_list/{handlers,services,repository,models,entities}` |
| Products — same layers | `src/product/{handlers,services,repository,models}` |
| WebSocket upgrade | `src/notification/handlers` |
| Config loading (env vars) | `config/` |
| MongoDB connection setup | `database/` |

## Getting started

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
