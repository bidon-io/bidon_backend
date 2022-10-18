class AppDemandProfile < ApplicationRecord
  belongs_to :account, polymorphic: true, class_name: 'DemandSourceAccount'
end
