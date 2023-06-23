# frozen_string_literal: true

class WinController < ApplicationController
  def create
    render_empty_result
  end

  private

  def schema_file_name
    'win.json'
  end
end
