# frozen_string_literal: true

class ConfigController < ApplicationController
  def create
    event = Events::Config.new(event_params)

    KafkaLogger.log(event)

    config_response = Api::Config::Response.new(api_request)

    if config_response.present?
      render json: config_response.body, status: :ok
    else
      render json:   { error: { code: 422, message: 'No adapters found' } },
             status: :unprocessable_entity
    end
  end

  private

  def schema_file_name
    'config.json'
  end
end
