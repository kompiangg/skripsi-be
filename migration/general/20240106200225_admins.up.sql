CREATE TABLE IF NOT EXISTS admins (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (account_id) REFERENCES accounts (id)
);