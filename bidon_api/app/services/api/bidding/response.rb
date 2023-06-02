# frozen_string_literal: true

module Api
  module Bidding
    DemandResponse = Struct.new(:demand, :raw_request, :raw_response, :status, :price, :seatbid, keyword_init: true)

    class Response
      prepend MemoWise

      attr_reader :bid_request

      delegate :app, :ad_type, :device, to: :bid_request

      def initialize(bid_request)
        @bid_request = bid_request
      end

      def bid?
        seat_bid[:price] > 0
      end

      def body
        {
          id:      SecureRandom.uuid,
          seatbid: seat_bid[:seatbid],
        }
      end
      memo_wise :body

      private

      def imp
        bid_request.params['imp'][0]
      end
      memo_wise :imp

      def seat_bid
        demands = imp.dig(:ext, :bidon, :bidding)

        responses = demands.map do |demand, hash|
          request_demand(demand, hash[:token], imp[:bidfloor]).call.tap do |demand_response|
            log(demand_response)
          end
        end

        responses.max_by(&:price)
      end
      memo_wise :seat_bid

      def request_demand(demand, token, bidfloor)
        case demand.to_s
        when 'bidmachine'
          Bidding::Demand::BidMachine.new(bid_request, token, bidfloor)
        else
          -> { DemandResponse.new(price: 0) }
        end
      end

      def log(demand_response)
        Rails.logger.info(demand_response.to_h.merge(action: 'bid').to_json)
      end
    end
  end
end
