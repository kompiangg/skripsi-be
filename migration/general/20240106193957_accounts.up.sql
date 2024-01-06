CREATE TABLE IF NOT EXISTS accounts (
  id UUID PRIMARY KEY,
  username VARCHAR(255),
  password VARCHAR(255),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);