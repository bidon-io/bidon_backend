# frozen_string_literal: true

class ShowController < ApplicationController
  def create
    event = Events::Show.new(event_params)

    KafkaLogger.log(event)

    render_empty_result
  end

  private

  def schema_file_name
    'show.json'
  end
end
