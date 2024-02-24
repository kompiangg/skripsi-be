CREATE TABLE IF NOT EXISTS orders (
  id varchar(26) PRIMARY KEY,
  cashier_id UUID,
  store_id varchar(10),
  payment_id varchar(10),
  customer_id varchar(10),
  total_quantity INT,
  total_unit VARCHAR(255),
  total_price DECIMAL(15,5),
  total_price_in_usd DECIMAL(15,5),
  currency VARCHAR(3),
  usd_rate DECIMAL(15,5),
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS order_details (
  id UUID PRIMARY KEY,
  order_id varchar(26),
  item_id varchar(10),
  quantity INT,
  unit VARCHAR(255),
  price DECIMAL(10,2),
  FOREIGN KEY (order_id) REFERENCES orders (id)
);