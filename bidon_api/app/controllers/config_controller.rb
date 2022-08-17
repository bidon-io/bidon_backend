# frozen_string_literal: true

class ConfigController < ApplicationController

  def create
    config_request = Config::Request.new(params)

    if config_request.valid?
      config_response = Config::Response.new(config_request)

      if config_response.present?
        render json: config_response.body, status: :ok
      else
        render_empty_result
      end
    else
      render_unprocessable_entity
    end
  end

  private

  def render_empty_result
    render json: { success: true }, status: :ok
  end

  def render_unprocessable_entity
    render json: { error: { code: 422, message: 'App key is invalid' } }, status: :unprocessable_entity
  end
end
