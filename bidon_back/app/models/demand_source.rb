class DemandSource < ApplicationRecord
  has_many :demand_source_accounts, dependent: :restrict_with_exception

  validates :api_key, presence: true, uniqueness: true
  validates :human_name, presence: true
end
