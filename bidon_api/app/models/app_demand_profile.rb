class AppDemandProfile < ApplicationRecord
  belongs_to :demand_source_account, key: :account_id
end
