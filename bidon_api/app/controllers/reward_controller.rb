# frozen_string_literal: true

class RewardController < ApplicationController
  def create
    event = Events::Reward.new(event_params)

    KafkaLogger.log(event)

    render_empty_result
  end

  private

  def schema_file_name
    'show.json'
  end
end
