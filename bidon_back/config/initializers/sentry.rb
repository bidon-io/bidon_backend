# frozen_string_literal: true

Sentry.init do |config|
  config.dsn = ENV.fetch('SENTRY_DSN')
  config.enabled_environments = %w[production staging]
  config.environment = Rails.env
  config.send_modules = false # include module versions in reports
  config.context_lines = 5 # number of lines of code context to capture

  config.send_default_pii = true
end

Sentry.set_tags(app: 'Bidon Backend')
