CREATE TABLE IF NOT EXISTS stores (
  id varchar(10) PRIMARY KEY,
  nation VARCHAR(10),
  region VARCHAR(255),
  district VARCHAR(255),
  sub_district VARCHAR(255),
  currency VARCHAR(3),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);