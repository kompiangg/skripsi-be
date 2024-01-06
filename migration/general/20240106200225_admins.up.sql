CREATE TABLE IF NOT EXISTS admins (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (account_id) REFERENCES accounts (id)
);