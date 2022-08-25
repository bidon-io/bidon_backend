# frozen_string_literal: true

class ConfigController < ApplicationController
  def create
    api_request = Api::Request.new(zipped_params.to_unsafe_h)

    log_params('CONFIG')

    if api_request.valid?
      config_response = Api::Config::Response.new(api_request)

      if config_response.present?
        render json: config_response.body, status: :ok
      else
        render_empty_result
      end
    else
      log_params('INVALID CONFIG')

      render_app_key_invalid
    end
  end
end
