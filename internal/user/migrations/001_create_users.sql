CREATE TABLE IF NOT EXISTS users (
    id CHAR(26) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    roles TEXT[] NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    seller_id CHAR(26),
    account_id CHAR(26)
);
