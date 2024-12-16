-- +goose Up
-- +goose StatementBegin
-- Drop the existing indexes
DROP INDEX IF EXISTS auction_configurations_default_uniq_idx;
DROP INDEX IF EXISTS auction_configurations_default_segment_uniq_idx;
DROP INDEX IF EXISTS index_app_demand_profiles_on_app_id_and_demand_source_id;
DROP INDEX IF EXISTS index_app_demand_profiles_on_public_uid;
DROP INDEX IF EXISTS index_apps_on_app_key;
DROP INDEX IF EXISTS index_apps_on_package_name_and_platform_id;
DROP INDEX IF EXISTS index_apps_on_public_uid;
DROP INDEX IF EXISTS index_auction_configurations_on_public_uid;
DROP INDEX IF EXISTS index_demand_source_accounts_on_public_uid;
DROP INDEX IF EXISTS index_line_items_on_public_uid;

-- Create updated indexes considering 'deleted_at'
CREATE UNIQUE INDEX auction_configurations_default_uniq_idx
    ON auction_configurations USING btree (app_id, ad_type, is_default)
    WHERE is_default IS TRUE AND deleted_at IS NULL;

CREATE UNIQUE INDEX auction_configurations_default_segment_uniq_idx
    ON auction_configurations USING btree (app_id, ad_type, is_default, segment_id)
    WHERE is_default IS FALSE AND deleted_at IS NULL;

CREATE UNIQUE INDEX index_app_demand_profiles_on_app_id_and_demand_source_id
    ON public.app_demand_profiles USING btree (app_id, demand_source_id)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX index_app_demand_profiles_on_public_uid
    ON public.app_demand_profiles USING btree (public_uid)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX index_apps_on_app_key
    ON public.apps USING btree (app_key)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX index_apps_on_package_name_and_platform_id
    ON public.apps USING btree (package_name, platform_id)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX index_apps_on_public_uid
    ON public.apps USING btree (public_uid)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX index_auction_configurations_on_public_uid
    ON public.auction_configurations USING btree (public_uid)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX index_demand_source_accounts_on_public_uid
    ON public.demand_source_accounts USING btree (public_uid)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX index_line_items_on_public_uid
    ON public.line_items USING btree (public_uid)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop updated indexes
DROP INDEX IF EXISTS auction_configurations_default_uniq_idx;
DROP INDEX IF EXISTS auction_configurations_default_segment_uniq_idx;
DROP INDEX IF EXISTS index_app_demand_profiles_on_app_id_and_demand_source_id;
DROP INDEX IF EXISTS index_app_demand_profiles_on_public_uid;
DROP INDEX IF EXISTS index_apps_on_app_key;
DROP INDEX IF EXISTS index_apps_on_package_name_and_platform_id;
DROP INDEX IF EXISTS index_apps_on_public_uid;
DROP INDEX IF EXISTS index_auction_configurations_on_public_uid;
DROP INDEX IF EXISTS index_demand_source_accounts_on_public_uid;
DROP INDEX IF EXISTS index_line_items_on_public_uid;

-- Recreate original indexes
CREATE UNIQUE INDEX auction_configurations_default_uniq_idx
    ON auction_configurations USING btree (app_id, ad_type, is_default)
    WHERE is_default IS TRUE;

CREATE UNIQUE INDEX auction_configurations_default_segment_uniq_idx
    ON auction_configurations USING btree (app_id, ad_type, is_default, segment_id)
    WHERE is_default IS FALSE;

CREATE UNIQUE INDEX index_app_demand_profiles_on_app_id_and_demand_source_id
    ON public.app_demand_profiles USING btree (app_id, demand_source_id);

CREATE UNIQUE INDEX index_app_demand_profiles_on_public_uid
    ON public.app_demand_profiles USING btree (public_uid);

CREATE UNIQUE INDEX index_apps_on_app_key
    ON public.apps USING btree (app_key);

CREATE UNIQUE INDEX index_apps_on_package_name_and_platform_id
    ON public.apps USING btree (package_name, platform_id);

CREATE UNIQUE INDEX index_apps_on_public_uid
    ON public.apps USING btree (public_uid);

CREATE UNIQUE INDEX index_auction_configurations_on_public_uid
    ON public.auction_configurations USING btree (public_uid);

CREATE UNIQUE INDEX index_demand_source_accounts_on_public_uid
    ON public.demand_source_accounts USING btree (public_uid);

CREATE UNIQUE INDEX index_line_items_on_public_uid
    ON public.line_items USING btree (public_uid);
-- +goose StatementEnd
