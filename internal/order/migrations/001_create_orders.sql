CREATE TABLE IF NOT EXISTS orders (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    receiver_id TEXT NOT NULL,
    account_id TEXT NOT NULL,
    seller_id TEXT NOT NULL,
    delivery_id TEXT NOT NULL,
    basket_id TEXT NOT NULL
);
