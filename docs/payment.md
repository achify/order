# Payment Service

This microservice records payments for orders. When a payment is accepted the related order status is automatically changed to `paid`.

## Environment
- `DB_DSN` – PostgreSQL connection string

## Endpoints
- `POST /payments` – register a payment
- `GET /payments/{id}` – get payment details

Run migrations before starting the service:
```bash
psql -h localhost -U postgres -d order -f internal/payment/migrations/001_create_payments.sql
```

Start the service:
```bash
go run ./cmd/payment
```
Logs are written to `payment.log`.
