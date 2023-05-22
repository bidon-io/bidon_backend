# frozen_string_literal: true

class RtbController < ApplicationController
  skip_before_action :validate_request_schema!
  skip_before_action :validate_app!

  # Action just for logging purposes currently.
  def create
    render_empty_result
  end
end
