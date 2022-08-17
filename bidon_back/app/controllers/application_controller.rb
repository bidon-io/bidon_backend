class ApplicationController < ActionController::Base
  before_action :set_sentry_context

  def set_sentry_context
    Sentry.set_extras(params: params.to_unsafe_h, session: session.to_hash)
  end
end
