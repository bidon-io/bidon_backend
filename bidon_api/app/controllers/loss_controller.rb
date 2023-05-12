# frozen_string_literal: true

class LossController < ApplicationController
  def create
    event = Events::Loss.new(event_params)

    KafkaLogger.log(event)

    render_empty_result
  end

  private

  def schema_file_name
    'loss.json'
  end
end
