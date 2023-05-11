# frozen_string_literal: true

class LossController < ApplicationController
  def create
    KafkaLogger.log_loss(kafka_event)

    render_empty_result
  end

  private

  def schema_file_name
    'loss.json'
  end
end
