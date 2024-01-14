CREATE TABLE IF NOT EXISTS stores (
  id varchar(10) PRIMARY KEY,
  region VARCHAR(255),
  district VARCHAR(255),
  sub_district VARCHAR(255),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);