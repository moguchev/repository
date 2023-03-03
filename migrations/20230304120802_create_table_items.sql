-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items_stocks (
    id bigint PRIMARY KEY,
    name text NOT NULL,
    price int4 NOT NULL,
    tags int4[] NOT NULL DEFAULT '{}',
    Length int4,
    Width int4,
    Height int4
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
