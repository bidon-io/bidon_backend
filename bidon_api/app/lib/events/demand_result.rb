# frozen_string_literal: true

class Events::DemandResult < Events::Stats
  def initialize(params, round_index:, demand_index:)
    super(params)

    @round_index = round_index
    @demand_index = demand_index
  end

  # Since this is originated from Stats event, it has all of the same params.
  # We just override some of the params from the parent event.
  # It probably should be a separate event with a separate schema.
  def to_hash
    super.merge!(
      'event_type'                => 'demand_result',
      'timestamp'                 => resolve_timestamp,
      'stats__result__status'     => demand['status'],
      'stats__result__winner_id'  => demand['id'],
      'stats__result__ad_unit_id' => demand['ad_unit_id'],
      'stats__result__ecpm'       => demand['ecpm'],
      'round_id'                  => round['id'],
      'pricefloor'                => round['pricefloor'],
    )
  end

  private

  def round
    @round ||= @params['stats']['rounds'][@round_index]
  end

  def demand
    @demand ||= round['demands'][@demand_index]
  end

  def resolve_timestamp
    demand_timestamp = demand['fill_finish_ts'].presence || demand['bid_finish_ts'].presence
    return @params.timestamp unless demand_timestamp

    demand_timestamp = demand_timestamp.to_f / 1000
    # We don't really care what the timestamp is,
    # as long as it's less than the timestamp of the stats event
    return @params.timestamp if demand_timestamp > @params.timestamp

    demand_timestamp
  end
end
