class App < ApplicationRecord
  belongs_to :user

  has_many :app_demand_profiles,    dependent: :restrict_with_exception
  has_many :line_items,             dependent: :restrict_with_exception
  has_many :auction_configurations, dependent: :destroy

  validates :package_name, :app_key,    presence: true, uniqueness: true
  validates :platform_id,  :human_name, presence: true
end
