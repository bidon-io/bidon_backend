# frozen_string_literal: true

module Api
  module Auction
    class Response
      attr_reader :auction_request

      delegate :present?, to: :body

      def initialize(auction_request)
        @auction_request = auction_request
      end

      def body
        @body ||= {
          'rounds'     => [],
          'line_items' => [],
          'token'      => '{}',
          'min_price'  => 0,
        }
      end
    end
  end
end
