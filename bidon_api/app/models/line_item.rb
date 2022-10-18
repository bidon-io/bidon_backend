class LineItem < ApplicationRecord
  belongs_to :account, polymorphic: true, class_name: 'DemandSourceAccount'

  enum ad_type: AdType::ENUM
end
