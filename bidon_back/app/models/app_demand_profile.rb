class AppDemandProfile < ApplicationRecord
  belongs_to :app
  belongs_to :account, class_name: 'DemandSourceAccount'
  belongs_to :demand_source

  validates :app_id, uniqueness: { scope: :demand_source_id }

  def slug
    "app_#{app_id}_#{DemandSourceType::ENUM.key(account_type).to_s.underscore}_#{account_id}"
  end
end
