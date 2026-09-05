-- +goose Up
CREATE TABLE payments (
    id          UUID PRIMARY KEY,
    order_id    UUID NOT NULL,
    amount      NUMERIC(12, 2) NOT NULL,
    status      TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_payment_status CHECK (status IN ('pending', 'succeeded', 'failed'))
);

CREATE INDEX idx_payments_order_id ON payments(order_id);

-- +goose Down
DROP TABLE payments;