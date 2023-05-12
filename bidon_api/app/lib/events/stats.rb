# frozen_string_literal: true

class Events::Stats < Events::Base
  KAFKA_TOPIC = ENV.fetch('KAFKA_STATS_TOPIC')

  def to_hash
    super.merge!(
      'event_type' => 'stats',
    )
  end

  def auction_events
    @params['stats']['rounds'].each_with_index.flat_map do |round, round_index|
      events = round['demands'].each_index.map do |demand_index|
        Events::DemandResult.new(@params, round_index:, demand_index:)
      end

      events << Events::RoundResult.new(@params, round_index:)
    end
  end
end
