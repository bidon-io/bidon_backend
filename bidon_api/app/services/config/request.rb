# frozen_string_literal: true

module Config
  class Request
    attr_reader :params

    def initialize(params)
      @params = params
    end

    # TODO: check if app is valid
    def valid?
      true
    end
  end
end
