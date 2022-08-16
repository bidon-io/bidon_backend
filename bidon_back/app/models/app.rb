class App < ApplicationRecord
  belongs_to :user

  validates :package_name, presence: true, uniqueness: true
  validates :app_key, presence: true, uniqueness: true
  validates :platform_id, presence: true
  validates :human_name, presence: true
end
