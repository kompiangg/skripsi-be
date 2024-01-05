CREATE TABLE IF NOT EXISTS orders (
  id UUID PRIMARY KEY,
  payment_id UUID NOT NULL,
  customer_id UUID NOT NULL,
  item_id UUID NOT NULL,
  store_id UUID NOT NULL,
  quantity INT NOT NULL,
  unit VARCHAR(255) NOT NULL,
  price BIGINT NOT NULL,
  total_price BIGINT NOT NULL,
  discount BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);