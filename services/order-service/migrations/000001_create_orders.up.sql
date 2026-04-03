-- Order status enum
CREATE TYPE order_status AS ENUM (
  'PENDING',
  'CONFIRMED',
  'PREPARING',
  'OUT_FOR_DELIVERY',
  'DELIVERED',
  'CANCELLED'
);

-- Orders table
CREATE TABLE orders (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        UUID NOT NULL,
  status         order_status NOT NULL DEFAULT 'PENDING',
  total_amount   BIGINT NOT NULL,
  currency       VARCHAR(3) NOT NULL DEFAULT 'INR',
  delivery_addr  TEXT NOT NULL,
  notes          TEXT,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Order items
CREATE TABLE order_items (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id     UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  item_name    VARCHAR(255) NOT NULL,
  quantity     INT NOT NULL CHECK (quantity > 0),
  unit_price   BIGINT NOT NULL,
  total_price  BIGINT NOT NULL
);

-- Indexes (performance)
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status  ON orders(status);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);