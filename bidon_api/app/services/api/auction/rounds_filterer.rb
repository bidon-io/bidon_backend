module Api
  module Auction
    class RoundsFilterer
      prepend MemoWise

      attr_reader :rounds, :adapters

      def initialize(rounds:, adapters:)
        @rounds = rounds
        @adapters = adapters
      end

      def fetch
        return [] unless adapters_names

        rounds.each_with_object([]) do |round, result|
          filtered_demands = round['demands']&.select { |demand| demand.in?(adapters_names) }
          return result if filtered_demands.blank?

          result << round.merge('demands' => filtered_demands)
        end
      end

      def adapters_names
        adapters.keys
      end
      memo_wise :adapters_names
    end
  end
end
