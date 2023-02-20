class LineItem < ApplicationRecord
  belongs_to :app
  belongs_to :account, class_name: 'DemandSourceAccount'

  validates :bid_floor, numericality: { greater_than_or_equal_to: 0 }
end
