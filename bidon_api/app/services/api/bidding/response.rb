# frozen_string_literal: true

module Api
  module Bidding
    DemandResponse = Struct.new(:demand, :raw_request, :raw_response, :status, :price, :bid, keyword_init: true)

    class Response
      prepend MemoWise

      attr_reader :bid_request

      delegate :app, :ad_type, :device, to: :bid_request

      def initialize(bid_request)
        @bid_request = bid_request
      end

      def bid?
        auction_result[:price] > 0
      end

      def body
        { bid: auction_result.bid }
      end
      memo_wise :body

      private

      def imp
        bid_request.params['imp']
      end
      memo_wise :imp

      def auction_result
        return empty_demand_response unless imp[:demands]

        responses = imp[:demands].map do |demand, hash|
          request_demand(demand, hash[:token], imp[:bidfloor]).call.tap do |demand_response|
            log(demand_response)
          end
        end

        responses.max_by(&:price)
      end
      memo_wise :auction_result

      def request_demand(demand, token, bidfloor)
        case demand.to_s
        when 'bidmachine'
          Bidding::Demand::BidMachine.new(bid_request, token, bidfloor)
        else
          -> { empty_demand_response }
        end
      end

      def log(demand_response)
        Rails.logger.info(demand_response.to_h.merge(action: 'bid').to_json)
      end

      def empty_demand_response
        DemandResponse.new(price: 0, bid: {})
      end
    end
  end
end
