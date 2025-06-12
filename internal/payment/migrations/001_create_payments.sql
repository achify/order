CREATE TABLE IF NOT EXISTS payments (
    id CHAR(26) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    order_id CHAR(26) NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL
);
