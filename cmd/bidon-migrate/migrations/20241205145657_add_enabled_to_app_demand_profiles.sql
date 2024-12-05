-- +goose Up
-- +goose StatementBegin
ALTER TABLE app_demand_profiles
ADD COLUMN enabled boolean DEFAULT true NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE app_demand_profiles
DROP COLUMN enabled;
-- +goose StatementEnd
