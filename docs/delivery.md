# Delivery Microservice

This service stores delivery records and synchronizes parcel machine
locations from Omniva once per day.

## Environment
- `DB_DSN` – PostgreSQL connection string

## Endpoints
- `GET /deliveries` – list all deliveries
- `GET /deliveries/{id}` – get delivery
- `POST /deliveries` – create delivery
- `PATCH /deliveries/{id}` – update delivery status
- `DELETE /deliveries/{id}` – delete delivery
- `GET /locations/{provider}` – list synced pickup locations

Tracking information can be obtained by creating a delivery with
`provider` set to `omniva` and providing a `tracking_code`.

## Maintenance
Run migrations before starting the service:

```bash
psql -h localhost -U postgres -d order -f internal/delivery/migrations/001_create_deliveries.sql
```

Start the microservice:

```bash
go run ./cmd/delivery
```

The service writes logs to `delivery.log` and synchronizes parcel
machine locations daily.
