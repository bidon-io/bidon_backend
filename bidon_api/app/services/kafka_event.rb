class KafkaEvent
  prepend MemoWise

  attr_reader :params, :ip

  def initialize(params:, ip:)
    @params = params
    @ip = ip
  end

  def build
    fill_geo_data!

    params
  end
  memo_wise :build

  private

  def fill_geo_data!
    params['geo']['ip'] = ip
    params['geo']['country'] = geo_data[:country_code]
    params['geo']['country_id'] = geo_data[:country_id]
  end

  def geo_data
    OfflineGeocoder.find_geo_data(ip)
  end
  memo_wise :geo_data
end
