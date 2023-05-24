# frozen_string_literal: true

class IapController < ApplicationController
  def create
    event = Events::Iap.new(event_params)

    KafkaLogger.log(event)

    render_empty_result
  end

  private

  def schema_file_name
    'iap.json'
  end
end
