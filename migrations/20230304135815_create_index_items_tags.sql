-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS items_tags_idx ON items USING GIN (tags);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS items_tags_idx;
-- +goose StatementEnd
