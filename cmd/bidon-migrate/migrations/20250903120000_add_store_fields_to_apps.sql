-- +goose Up
-- +goose StatementBegin
ALTER TABLE apps
ADD COLUMN store_id character varying,
ADD COLUMN store_url character varying,
ADD COLUMN categories text[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE apps
DROP COLUMN store_id,
DROP COLUMN store_url,
DROP COLUMN categories;
-- +goose StatementEnd
