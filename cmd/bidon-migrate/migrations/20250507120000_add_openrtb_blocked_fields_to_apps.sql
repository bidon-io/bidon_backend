-- +goose Up
-- +goose StatementBegin
ALTER TABLE apps
ADD COLUMN badv TEXT,
ADD COLUMN bcat TEXT,
ADD COLUMN bapp TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE apps
DROP COLUMN badv,
DROP COLUMN bcat,
DROP COLUMN bapp;
-- +goose StatementEnd

