-- +goose Up
-- +goose StatementBegin
ALTER TABLE line_items DROP COLUMN code;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE line_items ADD COLUMN code VARCHAR(255) NOT NULL;
-- +goose StatementEnd
