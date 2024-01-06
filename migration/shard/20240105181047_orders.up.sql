CREATE TABLE IF NOT EXISTS orders (
  id UUID PRIMARY KEY,
  payment_id UUID,
  customer_id VARCHAR(36),
  cashier_id UUID,
  item_id UUID,
  store_id UUID,
  quantity INT,
  unit VARCHAR(255),
  price BIGINT,
  total_price BIGINT,
  created_at TIMESTAMPTZ NOT NULL
);