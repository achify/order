CREATE TABLE IF NOT EXISTS baskets (
    id CHAR(26) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    account_id CHAR(26) NOT NULL,
    total_price BIGINT NOT NULL DEFAULT 0
);
