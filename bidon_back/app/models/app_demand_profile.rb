class AppDemandProfile < ApplicationRecord
  belongs_to :app
  belongs_to :account, polymorphic: true, class_name: 'DemandSourceAccount'
  belongs_to :demand_source
end
