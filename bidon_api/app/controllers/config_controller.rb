# frozen_string_literal: true

class ConfigController < ApplicationController
  def create
    config_response = Api::Config::Response.new(api_request)

    render json: config_response.body, status: :ok
  end
end
