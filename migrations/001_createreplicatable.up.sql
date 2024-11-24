CREATE TABLE replicas (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('active', 'disabled', 'inactive')),
    healthcheck_endpoint VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
