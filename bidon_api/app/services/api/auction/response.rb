# frozen_string_literal: true

module Api
  module Auction
    class Response
      prepend MemoWise

      attr_reader :auction_request

      delegate :present?, to: :body
      delegate :app, to: :auction_request

      def initialize(auction_request)
        @auction_request = auction_request
      end

      def body
        {
          'rounds'     => auction_configuration&.rounds || [],
          'line_items' => line_items,
          'token'      => '{}',
          'min_price'  => auction_configuration&.pricefloor || 0,
        }
      end
      memo_wise :body

      def auction_configuration
        AuctionConfiguration.where(app_id: app.id).order(Sequel.desc(:created_at)).first
      end
      memo_wise :auction_configuration

      def line_items
        LineItem.eager(demand_source_account: :demand_source).where(app_id: app.id).map do |line_item|
          {
            id:         line_item.demand_source_account.demand_source.api_key,
            pricefloor: line_item.bid_floor.to_f,
            ad_unit_id: line_item.code,
          }
        end
      end
      memo_wise :line_items
    end
  end
end
