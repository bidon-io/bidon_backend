-- +goose Up
-- +goose StatementBegin
INSERT INTO demand_sources (api_key, human_name, created_at, updated_at)
VALUES
('bidmachine', 'BidMachine', NOW(), NOW()),
('admob', 'AdMob', NOW(), NOW()),
('applovin', 'AppLovin', NOW(), NOW()),
('unityads', 'Unity Ads', NOW(), NOW()),
('meta', 'Meta', NOW(), NOW()),
('dtexchange', 'DT Exchange', NOW(), NOW()),
('amazon', 'Amazon', NOW(), NOW()),
('bigoads', 'Bigo Ads', NOW(), NOW()),
('chartboost', 'Chartboost', NOW(), NOW()),
('gam', 'Google Ad Manager', NOW(), NOW()),
('inmobi', 'InMobi', NOW(), NOW()),
('ironsource', 'ironSource', NOW(), NOW()),
('mintegral', 'Mintegral', NOW(), NOW()),
('mobilefuse', 'MobileFuse', NOW(), NOW()),
('moloco', 'Moloco', NOW(), NOW()),
('startio', 'Start.io', NOW(), NOW()),
('taurusx', 'TaurusX', NOW(), NOW()),
('vkads', 'VK Ads', NOW(), NOW()),
('vungle', 'Vungle', NOW(), NOW()),
('yandex', 'Yandex', NOW(), NOW())
ON CONFLICT (api_key) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
