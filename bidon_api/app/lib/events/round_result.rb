# frozen_string_literal: true

class Events::RoundResult < Events::Stats
  def initialize(params, round_index:)
    super(params)

    @round_index = round_index
  end

  # Since this is originated from Stats event, it has all of the same params.
  # We just override some of the params from the parent event.
  # It probably should be a separate event with a separate schema.
  def to_hash
    super.merge!(
      'event_type'                => 'round_result',
      'timestamp'                 => resolve_timestamp,
      'stats__result__status'     => round['winner_id'].present? ? 'SUCCESS' : 'FAIL',
      'stats__result__winner_id'  => round['winner_id'],
      'stats__result__ad_unit_id' => winner_demand['ad_unit_id'],
      'stats__result__ecpm'       => round['winner_ecpm'],
      'round_id'                  => round['id'],
      'pricefloor'                => round['pricefloor'],
    )
  end

  private

  def round
    @round ||= @params['stats']['rounds'][@round_index]
  end

  def winner_demand
    @winner_demand ||= round['demands'].find(-> { {} }) { _1['status'] == 'WIN' }
  end

  def resolve_timestamp
    round_timestamp = max_demands_timestamp
    return @params.timestamp unless round_timestamp

    round_timestamp = round_timestamp.to_f / 1000
    # We don't really care what the timestamp is,
    # as long as it's less than the timestamp of the StatsEvent
    return @params.timestamp if round_timestamp > @params.timestamp

    round_timestamp
  end

  def max_demands_timestamp
    ts = round['demands'].map do |demand|
      (demand['fill_finish_ts'].presence || demand['bid_finish_ts'].presence).to_i
    end.max

    return nil if ts == 0

    ts
  end
end
