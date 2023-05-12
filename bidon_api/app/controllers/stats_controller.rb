# frozen_string_literal: true

class StatsController < ApplicationController
  def create
    event = Events::Stats.new(event_params)
    auction_events = event.auction_events

    KafkaLogger.log_many([*auction_events, event])

    render_empty_result
  end

  private

  def schema_file_name
    'stats.json'
  end
end
