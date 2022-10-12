# frozen_string_literal: true

class StatsController < ApplicationController
  def create
    kafka_event = KafkaEvent.new(params: permitted_params, ip: request.remote_ip).build
    KafkaLogger.log_stats(kafka_event)

    render_empty_result
  end

  private

  def schema_file_name
    'stats.json'
  end
end
