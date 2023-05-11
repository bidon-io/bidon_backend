# frozen_string_literal: true

class ShowController < ApplicationController
  def create
    KafkaLogger.log_show(kafka_event)

    render_empty_result
  end

  private

  def schema_file_name
    'show.json'
  end
end
