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
        rounds.filter_map do |round|
          filtered_demands = round['demands'] & adapters_names
          next if filtered_demands.empty?

          round.merge('demands' => filtered_demands)
        end
      end

      def adapters_names
        adapters.keys
      end
      memo_wise :adapters_names
    end
  end
end
