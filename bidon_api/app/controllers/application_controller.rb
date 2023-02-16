class ApplicationController < ActionController::API
  prepend MemoWise

  before_action :set_sentry_context
  before_action :validate_bidon_header!
  before_action :validate_app!
  before_action :validate_request_schema!

  rescue_from StandardError, with: :handle_exception

  wrap_parameters false

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

  def validate_request_schema!
    return if schema_errors.none?

    render json:   { error: { code: 422, message: 'Invalid request schema', errors: schema_errors } },
           status: :unprocessable_entity
  end

  def handle_exception(error)
    Sentry.capture_exception(error)
    render json: { error: { code: 500, message: 'Internal Server Error' } }, status: :internal_server_error
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

  def remote_ip
    request.remote_ip
  end
  memo_wise :remote_ip

  def schema_errors
    schemer.validate(permitted_params).map do |error|
      error.slice('data_pointer', 'type', 'details')
    end
  end
  memo_wise :schema_errors

  def schemer
    Rails.cache.fetch("schemer_#{schema_file_name}") do
      JSONSchemer.schema(schema_path, ref_resolver: SchemerFileResolver.new)
    end
  end

  def schema_path
    Pathname.new(Rails.root.join('json_schema', schema_file_name))
  end

  def schema_file_name
    raise NotImplementedError
  end
end
