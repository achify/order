# Order Service

This service implements a simple but scalable order API using Domain Driven Design (DDD). Orders are stored in PostgreSQL and authenticated via JWT.

## Requirements

- Go 1.24+
- Docker (for Postgres via `docker-compose`)

## Running locally

Start Postgres:

```bash
docker-compose up -d
```

Run database migration:

```bash
psql -h localhost -U postgres -d order -f internal/order/migrations/001_create_orders.sql
```

Configure environment variables in a `.env` file (example values shown):

```bash
DB_DSN=postgres://postgres:postgres@localhost:5432/order?sslmode=disable
JWT_SECRET=secret
```

Start the API:

```bash
go run ./cmd/server
```

Logs are written to `server.log`. Tail them with:

```bash
tail -f server.log
```

## API

All endpoints require a valid JWT in the `Authorization: Bearer <token>` header.

- `GET /orders` – list orders
- `GET /orders/{id}` – get order by id
- `POST /orders` – create order
- `PATCH /orders/{id}` – update order
- `DELETE /orders/{id}` – delete order
- `GET /metrics` – basic memory and RPS metrics

Validation errors return HTTP `422`. Missing or invalid tokens return `401`. When an order isn't found `404` is used.
