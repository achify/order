CREATE TABLE IF NOT EXISTS categories (
    id CHAR(26) PRIMARY KEY,
    name TEXT NOT NULL,
    parent_id CHAR(26) REFERENCES categories(id)
);
