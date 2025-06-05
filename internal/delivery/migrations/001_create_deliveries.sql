CREATE TABLE IF NOT EXISTS deliveries (
    id CHAR(26) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    provider VARCHAR(50) NOT NULL,
    tracking_code VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS locations (
    id VARCHAR(50) PRIMARY KEY,
    provider VARCHAR(50) NOT NULL,
    data JSONB NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
