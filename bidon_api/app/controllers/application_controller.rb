class ApplicationController < ActionController::API
  prepend MemoWise

  before_action :set_sentry_context
  before_action :validate_bidon_header!
  before_action :validate_app!

  wrap_parameters false

  rescue_from StandardError do |error|
    Sentry.capture_exception(error)
    render json: { error: { code: 500, message: 'Internal Server Error' } }, status: :internal_server_error
  end

  private

  def validate_bidon_header!
    return if request.headers['X-BidOn-Version'].present?

    render json:   { error: { code: 422, message: 'Request should contain X-BidOn-Version header' } },
           status: :unprocessable_entity
  end

  def validate_app!
    return if api_request.valid?

    render json:   { error: { code: 422, message: 'App key is invalid' } },
           status: :unprocessable_entity
  end

  def render_empty_result
    render json: { success: true }, status: :ok
  end

  def set_sentry_context
    Sentry.set_extras(params:, session: session.to_hash)
  end

  def api_request
    Api::Request.new(permitted_params)
  end
  memo_wise :api_request

  def permitted_params
    params.except(:controller, :action).permit!.to_h
  end
  memo_wise :permitted_params
end
