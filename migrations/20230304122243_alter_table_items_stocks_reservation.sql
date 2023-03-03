-- +goose Up
-- +goose StatementBegin
ALTER TABLE items_stocks
    ADD COLUMN IF NOT EXISTS count int4 NOT NULL DEFAULT 0;

ALTER TABLE items_stocks_reservation
    ADD COLUMN IF NOT EXISTS count int4 NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE items_stocks_reservation
    DROP COLUMN IF EXISTS count;

ALTER TABLE items_stocks
    DROP COLUMN IF EXISTS count;
-- +goose StatementEnd
