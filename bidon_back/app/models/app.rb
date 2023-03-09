class App < ApplicationRecord
  belongs_to :user

  has_many :app_demand_profiles,    dependent: :restrict_with_exception
  has_many :line_items,             dependent: :restrict_with_exception
  has_many :auction_configurations, dependent: :destroy

  validates :app_key, presence: true, uniqueness: true
  validates :package_name, :platform_id, :human_name, presence: true
  validates :package_name, uniqueness: { scope: :platform_id }

  enum platform_id: { ios: 4, android: 1 }

  def settings=(value)
    super(JSON.parse(value.gsub('=>', ':')))
  end
end
