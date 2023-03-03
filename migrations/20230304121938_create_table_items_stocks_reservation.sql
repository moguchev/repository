-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items_stocks_reservation (
    items_id bigint,
    warehouse_id bigint,
    PRIMARY KEY (items_id, warehouse_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items_stocks_reservation;
-- +goose StatementEnd
