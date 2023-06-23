# frozen_string_literal: true

module Api
  module Auction
    class Response
      prepend MemoWise

      attr_reader :auction_request

      delegate :present?, to: :body
      delegate :app, :ad_type, :ad_object, :adapters, :device, to: :auction_request

      def initialize(auction_request)
        @auction_request = auction_request
      end

      def body
        return unless auction_configuration
        return if rounds.empty? # if we filtered all rounds return 422 response

        {
          'rounds'                     => rounds,
          'line_items'                 => line_items,
          'token'                      => '{}',
          'pricefloor'                 => ad_object['pricefloor'],
          'auction_id'                 => auction_id,
          'auction_configuration_id'   => auction_configuration.id,
          'external_win_notifications' => auction_configuration.external_win_notifications,
        }
      end
      memo_wise :body

      def rounds
        RoundsFilterer.new(rounds: JSON.parse(auction_configuration.rounds), adapters:).fetch
      end
      memo_wise :rounds

      def line_items
        LineItemsFetcher.new(
          app:, ad_type:, adapters:,
          banner_format: ad_object.dig('banner', 'format'),
          device_type:   device['type']
        ).fetch
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
