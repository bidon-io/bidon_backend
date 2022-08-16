class Country < ApplicationRecord
  validates :alpha_2_code, :alpha_3_code, uniqueness: true
end
