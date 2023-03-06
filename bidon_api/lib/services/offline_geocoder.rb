class OfflineGeocoder
  include Singleton
  prepend MemoWise

  DEFAULT_COUNTRY_CODES_FOR_CONTINENTS = {
    'Europe' => 'FR',
    'Asia'   => 'ID',
  }.freeze
  DEFAULT_COUNTRY_CODE = 'US'.freeze

  class << self
    delegate :find_geo_data, to: :instance
  end

  def find_geo_data(ip)
    country_code = country_code_for(ip)
    country_id = Country.find_cached(country_code)&.id

    { country_code:, country_id: }
  end

  def country_code_for(ip)
    geo_data = third_party_geo_data_for(ip)

    if geo_data.is_a?(MaxMindDB::Result)
      geo_data.country.iso_code || DEFAULT_COUNTRY_CODES_FOR_CONTINENTS[geo_data.continent.name] || DEFAULT_COUNTRY_CODE
    else
      geo_data.country[:iso] || DEFAULT_COUNTRY_CODE
    end
  end

  def third_party_geo_data_for(ip)
    max_mind = max_mind_db.lookup(ip)

    return max_mind if max_mind&.city

    sypex_db.query(ip)
  end

  def max_mind_db
    MaxMindDB.new(Utils.fetch_from_env('MAXMIND_GEOIP_FILE_PATH'))
  end
  memo_wise :max_mind_db

  def sypex_db
    SypexGeo::Database.new(Utils.fetch_from_env('SYPEX_GEOIP_FILE_PATH'))
  end
  memo_wise :sypex_db
end
