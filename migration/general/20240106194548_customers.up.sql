CREATE TABLE IF NOT EXISTS customers (
  id varchar(36) PRIMARY KEY,
  account_id UUID NOT NULL,
  name VARCHAR(255) NOT NULL,
  contact VARCHAR(50) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  FOREIGN KEY (account_id) REFERENCES accounts (id)
);