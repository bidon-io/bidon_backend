# frozen_string_literal: true

class RewardController < ApplicationController
  def create
    KafkaLogger.log_reward(kafka_event)

    render_empty_result
  end

  private

  def schema_file_name
    'show.json'
  end
end
