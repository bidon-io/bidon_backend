-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.auction_configurations
    ADD COLUMN timeout integer NOT NULL DEFAULT 0,
    ADD COLUMN demands text[] DEFAULT ARRAY[]::text[],
    ADD COLUMN biddings text[] DEFAULT ARRAY[]::text[],
    ADD COLUMN ad_unit_ids bigint[] DEFAULT ARRAY[]::bigint[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.auction_configurations
    DROP COLUMN timeout,
    DROP COLUMN demands,
    DROP COLUMN biddings,
    DROP COLUMN ad_unit_ids;
-- +goose StatementEnd
