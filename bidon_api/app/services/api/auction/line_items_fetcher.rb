module Api
  module Auction
    class LineItemsFetcher
      prepend MemoWise

      FORMATS = %w[BANNER LEADERBOARD MREC ADAPTIVE].freeze

      attr_reader :app, :ad_type, :adapters, :banner_format

      def initialize(app:, ad_type:, adapters:, banner_format:)
        @app = app
        @ad_type = ad_type
        @adapters = adapters
        @banner_format = banner_format
      end

      def fetch
        line_items.filter_map do |line_item|
          api_key = line_item.demand_source_account.demand_source.api_key
          next unless adapters.key?(api_key)

          {
            id:         api_key,
            pricefloor: line_item.bid_floor.to_f,
            ad_unit_id: line_item.code,
          }
        end
      end

      private

      def line_items
        result = LineItem.eager(demand_source_account: :demand_source)
                         .where(app_id: app.id, ad_type: AdType::ENUM[ad_type])

        return result unless ad_type == :banner
        return [] unless FORMATS.include?(banner_format)

        result.where(format: banner_format)
      end
    end
  end
end
