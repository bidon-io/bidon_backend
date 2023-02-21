# frozen_string_literal: true

module Api
  module Auction
    class Response
      prepend MemoWise

      attr_reader :auction_request

      delegate :present?, to: :body
      delegate :app, :ad_type, :ad_object, to: :auction_request

      def initialize(auction_request)
        @auction_request = auction_request
      end

      def body
        return unless auction_configuration

        {
          'rounds'                   => rounds,
          'line_items'               => line_items,
          'token'                    => '{}',
          'pricefloor'               => auction_configuration.pricefloor,
          'auction_id'               => auction_id,
          'auction_configuration_id' => auction_configuration.id,
        }
      end
      memo_wise :body

      def rounds
        JSON.parse(auction_configuration.rounds)
      end
      memo_wise :rounds

      def line_items
        LineItemsFetcher.new(app:, ad_type:, banner_format: ad_object.dig('banner', 'format')).fetch
      end
      memo_wise :line_items

      def auction_id
        SecureRandom.uuid
      end
      memo_wise :auction_id

      def auction_configuration
        AuctionConfiguration.where(app_id: app.id, ad_type: AdType::ENUM[ad_type])
                            .order(Sequel.desc(:created_at)).first
      end
      memo_wise :auction_configuration
    end
  end
end
