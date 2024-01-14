CREATE TABLE IF NOT EXISTS orders (
  id UUID PRIMARY KEY,
  cashier_id UUID,
  item_id varchar(10),
  store_id varchar(10),
  payment_id varchar(10),
  customer_id varchar(10),
  quantity INT,
  unit VARCHAR(255),
  price DECIMAL(10,2),
  total_price DECIMAL(10,2),
  created_at TIMESTAMPTZ NOT NULL
);