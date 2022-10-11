# frozen_string_literal: true

class ClickController < ApplicationController
  def create
    render_empty_result
  end

  private

  def schema_file_name
    'show.json'
  end
end
