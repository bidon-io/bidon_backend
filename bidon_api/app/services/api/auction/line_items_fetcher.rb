module Api
  module Auction
    class LineItemsFetcher
      prepend MemoWise

      # rubocop:disable Style/MutableConstant
      BANNER_FORMAT      = 'BANNER'
      LEADERBOARD_FORMAT = 'LEADERBOARD'
      MREC_FORMAT        = 'MREC'
      ADAPTIVE_FORMAT    = 'ADAPTIVE'
      # rubocop:enable Style/MutableConstant

      FORMAT_SIZES = {
        BANNER_FORMAT      => { width: 320, height: 50 },
        LEADERBOARD_FORMAT => { width: 728, height: 90 },
        MREC_FORMAT        => { width: 300, height: 250 },
        ADAPTIVE_FORMAT    => { width: 0,   height: 50 },
      }.freeze

      attr_reader :app, :ad_type, :banner_format

      def initialize(app:, ad_type:, banner_format: 0)
        @app = app
        @ad_type = ad_type
        @banner_format = banner_format
      end

      def fetch
        line_items.map do |line_item|
          {
            id:         line_item.demand_source_account.demand_source.api_key,
            pricefloor: line_item.bid_floor.to_f,
            ad_unit_id: line_item.code,
          }
        end
      end

      private

      def line_items
        result = LineItem.eager(demand_source_account: :demand_source)
                         .where(app_id: app.id, ad_type: AdType::ENUM[ad_type])

        if ad_type == :banner
          result.where(FORMAT_SIZES[banner_format])
        else
          result
        end
      end
    end
  end
end
