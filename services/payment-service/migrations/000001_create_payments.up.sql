-- Payment status enum
CREATE TYPE payment_status AS ENUM (
  'CREATED',
  'PROCESSING',
  'SUCCESS',
  'FAILED',
  'INVALID'
);

CREATE TABLE payments (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id         UUID NOT NULL,
  idempotency_key  VARCHAR(255) NOT NULL,
  amount           BIGINT NOT NULL,
  currency         VARCHAR(3) NOT NULL DEFAULT 'INR',
  status           payment_status NOT NULL DEFAULT 'CREATED',
  gateway_ref      VARCHAR(255),
  failure_reason   TEXT,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT uq_idempotency_key UNIQUE (idempotency_key)
);

-- Indexes
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_status   ON payments(status);