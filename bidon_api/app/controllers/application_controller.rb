class ApplicationController < ActionController::API
  before_action :set_sentry_context
  before_action :validate_bidon_header!

  rescue_from StandardError do |_e|
    render json: { error: { code: 500, message: 'Internal Server Error' } }, status: :internal_server_error
  end

  private

  def validate_bidon_header!
    return if request.headers['X-BidOn-Version'].present?

    render json:   { error: { code: 422, message: 'Request should contain X-BidOn-Version header' } },
           status: :unprocessable_entity
  end

  def render_empty_result
    render json: { success: true }, status: :ok
  end

  def render_app_key_invalid
    render json: { error: { code: 422, message: 'App key is invalid' } }, status: :unprocessable_entity
  end

  def zipped_params
    json = Utils.decode_params(request.raw_post)
    ActionController::Parameters.new(JSON.parse(json))
  rescue Zlib::GzipFile::Error, JSON::ParserError
    ActionController::Parameters.new
  end

  def set_sentry_context
    Sentry.set_extras(params: params.to_unsafe_h, session: session.to_hash)
  end
end
