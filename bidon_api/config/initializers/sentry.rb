# frozen_string_literal: true

Sentry.init do |config|
  config.dsn = ENV.fetch('SENTRY_DSN')
  config.enabled_environments = %w[production staging]
  config.environment = Rails.env
  config.send_modules = false
  config.context_lines = 5

  config.send_default_pii = true
end

Sentry.set_tags(app: 'Bidon SDK API')
