class OfflineGeocoder
  include Singleton
  prepend MemoWise

  # https://github.com/InteractiveAdvertisingBureau/AdCOM/blob/master/AdCOM%20v1.0%20FINAL.md#list--ip-location-services-
  MAX_MIND_PROVIDER_CODE = 3
  DEFAULT_COUNTRY_CODES_FOR_CONTINENTS = {
    'Europe' => 'FR',
    'Asia'   => 'ID',
  }.freeze

  Result = Struct.new(
    :country_code, :country_id, :country_code3, :city_name, :region_name, :region_code,
    :lat, :lon, :accuracy, :zip_code, :ip_service, :unknown_country,
    keyword_init: true
  )

  class << self
    delegate :find_geo_data, to: :instance
  end

  # @return [OfflineGeocoder::Result]
  def find_geo_data(ip) # rubocop:disable Metrics/AbcSize
    geo_data = lookup_ip(ip)
    country_code = country_code_for(geo_data)
    country = Country.find_cached(country_code)

    Result.new(
      country_code:,
      country_code3:   country&.alpha_3_code || Country::UNKNOWN_COUNTRY_CODE3,
      unknown_country: country_code == Country::UNKNOWN_COUNTRY_CODE,
      country_id:      country&.id,
      city_name:       geo_data.city.name,
      region_name:     geo_data.subdivisions.most_specific.name,
      region_code:     geo_data.subdivisions.most_specific.iso_code,
      lat:             geo_data.location.latitude,
      lon:             geo_data.location.longitude,
      accuracy:        geo_data.location.accuracy_radius.to_i * 1000, # convert kilometers to meters
      zip_code:        geo_data.postal.code,
      ip_service:      MAX_MIND_PROVIDER_CODE,
    )
  end

  def country_code_for(geo_data)
    geo_data.country.iso_code \
      || DEFAULT_COUNTRY_CODES_FOR_CONTINENTS[geo_data.continent.name] \
      || Country::UNKNOWN_COUNTRY_CODE
  end

  # @return [MaxMindDB::Result]
  def lookup_ip(ip)
    max_mind_db.lookup(ip)
  end

  def max_mind_db
    MaxMindDB.new(Utils.fetch_from_env('MAXMIND_GEOIP_FILE_PATH'))
  end
  memo_wise :max_mind_db
end
