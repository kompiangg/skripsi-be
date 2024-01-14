CREATE TABLE IF NOT EXISTS payment_types (
    id varchar(10) PRIMARY KEY,
    "type" VARCHAR(255) NOT NULL,
    bank VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);