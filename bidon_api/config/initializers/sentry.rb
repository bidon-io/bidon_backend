# frozen_string_literal: true

Sentry.init do |config|
  config.dsn = ENV.fetch('SENTRY_DSN', nil)
  config.enabled_environments = %w[production staging]
  config.environment = Rails.env
  config.send_modules = false
  config.context_lines = 5
  config.breadcrumbs_logger = %i[active_support_logger http_logger]
  config.send_default_pii = true

  # Performance monitoring
  config.traces_sample_rate = 1.0
end

Sentry.set_tags(app: 'Bidon SDK API')
