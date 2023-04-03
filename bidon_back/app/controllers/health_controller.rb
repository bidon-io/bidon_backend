# frozen_string_literal: true

# Default Rails implementation of a health check endpoint
# https://github.com/rails/rails/blob/main/railties/lib/rails/health_controller.rb
if Rails.gem_version >= Gem::Version.new('7.1')
  Rails.logger.warn 'Rails 7.1+ has a built-in health check endpoint, please remove this file'
end

class HealthController < ActionController::Base # rubocop:disable Rails/ApplicationController
  rescue_from(Exception) { render_down }

  def show
    render_up
  end

  private

  def render_up
    render html: html_status(color: 'green')
  end

  def render_down
    render html: html_status(color: 'red'), status: :internal_server_error
  end

  def html_status(color:)
    %(<html><body style="background-color: #{color}"></body></html>).html_safe # rubocop:disable Rails/OutputSafety
  end
end
