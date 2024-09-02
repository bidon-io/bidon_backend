-- +goose Up
-- +goose StatementBegin
ALTER TABLE auction_configurations
    ADD COLUMN is_default boolean DEFAULT false NOT NULL;

CREATE UNIQUE INDEX auction_configurations_default_uniq_idx ON auction_configurations (app_id, ad_type, is_default)
    WHERE is_default IS TRUE;
CREATE UNIQUE INDEX auction_configurations_default_segment_uniq_idx ON auction_configurations (app_id, ad_type, is_default, segment_id)
    WHERE is_default IS FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX auction_configurations_default_uniq_idx;
DROP INDEX auction_configurations_default_segment_uniq_idx;

ALTER TABLE auction_configurations
    DROP COLUMN is_default;
-- +goose StatementEnd
