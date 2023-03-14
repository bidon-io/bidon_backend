# frozen_string_literal: true

module Api
  class Request
    prepend MemoWise

    attr_reader :params

    def initialize(params)
      @params = params
    end

    def valid?
      app.present?
    end

    def ad_type
      params['ad_type'].to_sym
    end

    def ad_object
      params['ad_object']
    end

    def adapters
      params['adapters'].presence || {}
    end

    def app
      app_key = params.dig('app', 'key')
      package_name = params.dig('app', 'bundle')

      return unless app_key && package_name

      App.find(app_key:, package_name:)
    end
    memo_wise :app
  end
end
