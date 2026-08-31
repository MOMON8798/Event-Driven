-- +goose Up
CREATE TABLE orders (
    id          UUID PRIMARY KEY,
    client_id   TEXT NOT NULL,
    name        TEXT NOT NULL,
    total       NUMERIC(12, 2) NOT NULL,
    status      TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_client_id ON orders(client_id);

-- +goose Down
DROP TABLE orders;