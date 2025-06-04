CREATE TABLE IF NOT EXISTS orders (
    id CHAR(26) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    receiver_id CHAR(26) NOT NULL,
    account_id CHAR(26) NOT NULL,
    seller_id CHAR(26) NOT NULL,
    delivery_id CHAR(26) NOT NULL,
    basket_id CHAR(26) NOT NULL
);
