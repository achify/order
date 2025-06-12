# Order Service

This service implements a simple but scalable order API using Domain Driven Design (DDD). Orders are stored in PostgreSQL and authenticated via JWT.

OpenAPI specifications for all services are located in the `docs` directory.

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

# create users table
psql -h localhost -U postgres -d order -f internal/user/migrations/001_create_users.sql
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

### Delivery Service

The repository also includes a delivery microservice responsible for tracking
shipments and synchronizing parcel machine locations from Omniva. Run database
migrations located in `internal/delivery/migrations` and start the service:

```bash
psql -h localhost -U postgres -d order -f internal/delivery/migrations/001_create_deliveries.sql
go run ./cmd/delivery
```

Parcel machine data is fetched from `https://www.omniva.ee/locations.json` once a
day automatically. Logs are written to `delivery.log`.

API documentation is available in `docs/delivery-swagger.yaml`.

The main order API is documented in `docs/order-swagger.yaml`.

### Payment Service

Payments are recorded through a separate microservice. Run the migration and start the service:

```bash
psql -h localhost -U postgres -d order -f internal/payment/migrations/001_create_payments.sql
go run ./cmd/payment
```

See `docs/payment.md` for more details. Successful payments automatically set the related order status to `paid` and logs are written to `payment.log`.
API documentation for the payment service is available in `docs/payment-swagger.yaml`.

## API

All endpoints require a valid JWT in the `Authorization: Bearer <token>` header.

- `GET /orders` – list orders
- `GET /orders/{id}` – get order by id
- `POST /orders` – create order
- `PATCH /orders/{id}` – update order
- `DELETE /orders/{id}` – delete order
- `GET /metrics` – basic memory and RPS metrics

Validation errors return HTTP `422`. Missing or invalid tokens return `401`. When an order isn't found `404` is used.

### Authentication

Obtain a JWT token using credentials stored in the `users` table (create users via `/users`):

```bash
curl -X POST http://localhost:8080/auth/login \
  -d '{"username":"<user>","password":"<pass>"}' \
  -H 'Content-Type: application/json'
```

The response contains `token` and `refresh_token` fields. Tokens expire after 30 minutes by default. Refresh an access token with:

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -d '{"refresh_token":"<refresh>"}' \
  -H 'Content-Type: application/json'
```

### Example Order Request

```bash
# create order
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer <token>" \
  -H 'Content-Type: application/json' \
  -d '{"receiver_id":"r","account_id":"a","seller_id":"s","delivery_id":"d","basket_id":"b"}'
```
