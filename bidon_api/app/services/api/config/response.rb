# frozen_string_literal: true

module Api
  module Config
    class Response
      prepend MemoWise

      attr_reader :config_request

      delegate :present?, to: :body

      def initialize(config_request)
        @config_request = config_request
      end

      def body
        {
          'init'       => {
            'tmax'     => 5000,
            'adapters' => adapters,
          },
          'placements' => [],
          'token'      => '{}',
          'segment_id' => '',
        }
      end
      memo_wise :body

      def adapters
        AdaptersFetcher.new(app: config_request.app, config_adapters: config_request.adapters).fetch
      end
    end
  end
end
