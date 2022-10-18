class Country < ApplicationRecord
  def self.find_cached(code)
    Rails.cache.fetch("country_#{code}") do
      Country.find_by(alpha_2_code: code) || Country.find_by(alpha_2_code: 'ZZ')
    end
  end
end
