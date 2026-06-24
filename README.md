# Merkadapp API

REST API for a personal grocery management and expense tracking app. Built as a hands-on project to learn Go, following a layered architecture.

Pairs with [merkadapp_frontend](https://github.com/raulito1500/merkadapp_frontend) — the React SPA that consumes this API.

## Stack

- **Go 1.22**
- **Gin** — HTTP router and middleware
- **MongoDB** — via the official Go driver
- **Gorilla WebSocket** — real-time notifications channel
- **gin-contrib/cors** — CORS configuration

## Architecture

The codebase follows a strict handlers → services → repositories separation. Each domain owns its full stack independently:

```
src/
├── bill/
│   ├── entities/     # Response types (BillTotal, BillItem, Recommendation)
│   ├── handlers/     # HTTP layer — parses requests, calls service, writes responses
│   ├── models/       # MongoDB document models (Bill, BillItem, BillBags, BillTaxes)
│   ├── repository/   # MongoDB queries (interface + implementation)
│   └── services/     # Business logic (totals, merging, recommendations)
├── market_list/
│   ├── entities/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   └── services/
├── notification/
│   └── handlers/     # WebSocket upgrade and broadcast
└── product/
    ├── handlers/
    ├── models/
    ├── repository/
    └── services/
```

## API Reference

### Products

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/products/:id` | List products |
| `POST` | `/products` | Create product |
| `PUT` | `/products/:id` | Update product |
| `GET` | `/products/recommendations` | Products due for restocking based on purchase frequency |
| `GET` | `/products/:id/bill-items` | Full purchase history for a product |

### Bills

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/bills` | List all bills |
| `GET` | `/bills/:id` | Get bill by ID |
| `POST` | `/bills` | Create bill |
| `PUT` | `/bills/:id` | Update bill |
| `PUT` | `/bills/merge/:idDestination` | Merge multiple bills into one |
| `GET` | `/bills/byMonth` | Monthly spending totals for the last 6 months |

### Market Lists

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/market-list` | List all market lists |
| `GET` | `/market-list/:id` | Get market list by ID |
| `POST` | `/market-list` | Create market list |
| `GET` | `/market-list/suggested` | Auto-generate a list from product purchase frequency |
| `PUT` | `/market-list/:id/check/:idItem` | Mark a list item as purchased |

### Real-time

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/ws` | WebSocket connection for live push notifications |

## Local Setup

**Prerequisites:** Go 1.22+ and a MongoDB instance (local or [Atlas free tier](https://www.mongodb.com/cloud/atlas))

```bash
git clone https://github.com/raulito1500/merkadapp.git
cd merkadapp

cp .env.example .env
# edit .env and fill in your MongoDB credentials

go mod download
go run .
```

The server starts at `http://localhost:8080`.

### Environment variables

```env
DATABASE_URL=mongodb+srv://<user>:<password>@<cluster>.mongodb.net/
DATABASE_NAME=merkadapp
PORT=8080
```
