# frozen_string_literal: true

class LossController < ApplicationController
  def create
    Rails.logger.info("LossController#create: #{permitted_params}")
    render_empty_result
  end

  private

  def schema_file_name
    'loss.json'
  end
end
