# frozen_string_literal: true

class ConfigController < ApplicationController
  def create
    config_response = Api::Config::Response.new(api_request)

    if config_response.present?
      render json: config_response.body, status: :ok
    else
      render_empty_result
    end
  end
end
