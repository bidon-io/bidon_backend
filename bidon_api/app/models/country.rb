class Country < Sequel::Model
  UNKNOWN_COUNTRY_CODE = 'ZZ'.freeze
  UNKNOWN_COUNTRY_CODE3 = 'ZZZ'.freeze
  # @param [String] code ALPHA2 country code
  # @return [Country]
  def self.find_cached(code)
    Rails.cache.fetch("country_#{code}") do
      Country.find(alpha_2_code: code) || Country.find(alpha_2_code: UNKNOWN_COUNTRY_CODE)
    end
  end
end
