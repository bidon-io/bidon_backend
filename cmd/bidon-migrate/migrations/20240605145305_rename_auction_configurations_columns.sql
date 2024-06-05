-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.auction_configurations RENAME COLUMN biddings TO bidding;
ALTER TABLE public.auction_configurations ALTER COLUMN bidding TYPE character varying[] USING bidding::character varying[];
ALTER TABLE public.auction_configurations ALTER COLUMN demands TYPE character varying[] USING demands::character varying[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.auction_configurations RENAME COLUMN bidding TO biddings;
ALTER TABLE public.auction_configurations ALTER COLUMN biddings TYPE text[] USING biddings::text[];
ALTER TABLE public.auction_configurations ALTER COLUMN demands TYPE text[] USING demands::text[];
-- +goose StatementEnd
