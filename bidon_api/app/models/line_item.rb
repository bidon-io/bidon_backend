class LineItem < ApplicationRecord
  belongs_to :demand_source_account, key: :account_id

  enum ad_type: AdType::ENUM
end
