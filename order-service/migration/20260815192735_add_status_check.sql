-- +goose Up
ALTER TABLE orders ADD CONSTRAINT chk_status
    CHECK (status IN ('created','paid','shipped','delivered','cancelled'));

-- +goose Down
ALTER TABLE orders DROP CONSTRAINT chk_status;