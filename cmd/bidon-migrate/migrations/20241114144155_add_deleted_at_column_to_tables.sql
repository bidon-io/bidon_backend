-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.apps ADD COLUMN deleted_at timestamp(6) without time zone;
ALTER TABLE public.app_demand_profiles ADD COLUMN deleted_at timestamp(6) without time zone;
ALTER TABLE public.auction_configurations ADD COLUMN deleted_at timestamp(6) without time zone;
ALTER TABLE public.demand_source_accounts ADD COLUMN deleted_at timestamp(6) without time zone;
ALTER TABLE public.line_items ADD COLUMN deleted_at timestamp(6) without time zone;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.apps DROP COLUMN deleted_at;
ALTER TABLE public.app_demand_profiles DROP COLUMN deleted_at;
ALTER TABLE public.auction_configurations DROP COLUMN deleted_at;
ALTER TABLE public.demand_source_accounts DROP COLUMN deleted_at;
ALTER TABLE public.line_items DROP COLUMN deleted_at;
-- +goose StatementEnd
