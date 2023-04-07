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
    parse_ext!

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

  def parse_ext!
    return if params['ext'].blank?

    params['ext'] = JSON.parse(params['ext'])
  rescue JSON::ParserError => e
    Rails.logger.error("Failed to parse 'ext': #{e.message}")
    Sentry.capture_exception(e)
  end

  def geo_data
    OfflineGeocoder.find_geo_data(ip)
  end
  memo_wise :geo_data
end
