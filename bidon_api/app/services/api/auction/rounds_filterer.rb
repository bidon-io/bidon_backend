module Api
  module Auction
    class RoundsFilterer
      prepend MemoWise

      attr_reader :rounds, :adapters

      def initialize(rounds:, adapters:)
        @rounds = rounds
        @adapters = adapters
      end

      # @return [Array]
      def fetch
        rounds.filter_map do |round|
          demands = Array(round['demands']) & adapters_names
          bidding = Array(round['bidding']) & adapters_names

          next if demands.empty? && bidding.empty? # remove round if all empty

          round.merge('demands' => demands, 'bidding' => bidding)
        end
      end

      def adapters_names
        adapters.keys
      end
      memo_wise :adapters_names
    end
  end
end
