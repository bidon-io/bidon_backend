class Country < Sequel::Model
  def self.find_cached(code)
    Rails.cache.fetch("country_#{code}") do
      Country.find(alpha_2_code: code) || Country.find(alpha_2_code: 'ZZ')
    end
  end
end
