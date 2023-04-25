class KafkaEvent
  prepend MemoWise

  attr_reader :params, :ip

  def initialize(params:, ip:)
    @params = params
    @ip = ip
  end

  def build
    fill_timestamp!
    fill_geo_data!

    params['ext'] = parse_ext(params['ext'])
    fill_ext_with_empty_values_if_needed!

    alias_bid_to_show!

    params
  end
  memo_wise :build

  private

  def fill_timestamp!
    params['timestamp'] = Time.current.to_f
  end

  def fill_geo_data!
    params['geo'] ||= {}
    params['geo']['ip'] = ip
    params['geo']['country'] = geo_data[:country_code]
    params['geo']['country_id'] = geo_data[:country_id]
  end

  def alias_bid_to_show!
    return if params.key?('show')

    params['show'] = params['bid']
  end

  def parse_ext(ext)
    ext = parse_json(ext)
    # Android SDK version 2.6.40 sends double escaped JSON
    # TODO: Remove after SDK fixes this
    return ext if !ext.key?('appodeal_token') || ext['appodeal_token'].is_a?(Hash)

    ext['appodeal_token'] = parse_json(ext['appodeal_token'])

    ext
  rescue JSON::ParserError => e
    Rails.logger.error("Failed to parse 'ext': #{e.message}")
    Sentry.capture_exception(e)
  end

  def parse_json(source)
    return {} if source.blank?

    JSON.parse(source)
  end

  def fill_ext_with_empty_values_if_needed! # rubocop:disable Metrics/AbcSize
    params['ext']['appodeal_session_id'] ||= ''
    params['ext']['appodeal_segment_id'] ||= 0
    params['ext']['appodeal_placement_id'] ||= 0
    params['ext']['appodeal_token'] ||= {}
    params['ext']['appodeal_token']['signature'] ||= ''
  end

  def geo_data
    OfflineGeocoder.find_geo_data(ip)
  end
  memo_wise :geo_data
end
