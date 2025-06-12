CREATE TABLE IF NOT EXISTS items (
    id CHAR(26) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    price BIGINT NOT NULL,
    category_id CHAR(26) NOT NULL REFERENCES categories(id)
);
