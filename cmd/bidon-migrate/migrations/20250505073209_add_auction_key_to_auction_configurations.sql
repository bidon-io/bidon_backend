-- +goose Up
-- +goose StatementBegin
ALTER TABLE auction_configurations ADD COLUMN auction_key TEXT;
CREATE INDEX idx_auction_configurations_auction_key ON auction_configurations(auction_key);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_auction_configurations_auction_key;
ALTER TABLE auction_configurations DROP COLUMN IF EXISTS auction_key;
-- +goose StatementEnd
