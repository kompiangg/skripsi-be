CREATE TABLE IF NOT EXISTS cashiers (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL,
    store_id varchar(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (store_id) REFERENCES stores (id),
    FOREIGN KEY (account_id) REFERENCES accounts (id)
);