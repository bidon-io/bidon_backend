# frozen_string_literal: true

# EventParams wraps SDK request params and provides additional context
# to the data that all Events must have.
class EventParams
  attr_reader :timestamp

  delegate :[], to: :@request_params

  def initialize(request_params:, ip:)
    @request_params = request_params
    @ip = ip

    @timestamp = Time.current.to_f
  end

  def to_hash
    hash = @request_params.deep_dup

    hash['timestamp'] = timestamp
    hash['geo'] = geo
    hash['ext'] = ext

    hash
  end

  def geo
    @geo ||= build_geo
  end

  def ext
    @ext ||= parse_ext
  end

  private

  def build_geo
    geo_data = OfflineGeocoder.find_geo_data(@ip)

    geo = {
      'ip'         => @ip,
      'country'    => geo_data[:country_code],
      'country_id' => geo_data[:country_id],
    }

    return @request_params['geo'].merge(geo) if @request_params['geo'].present?

    geo
  end

  def parse_ext
    ext = parse_json(@request_params['ext'])
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
end
