# frozen_string_literal: true

class StatsController < ApplicationController
  def create
    render_empty_result
  end

  private

  def schema_file_name
    'stats.json'
  end
end
