class LineItem < ApplicationRecord
  belongs_to :app
  belongs_to :account, class_name: 'DemandSourceAccount'

  enum ad_type: AdType::ENUM

  validates :bid_floor, numericality: { greater_than_or_equal_to: 0 }

  def extra=(value)
    if value.is_a?(Hash)
      super(value)
    else
      super(JSON.parse(value.gsub('=>', ':')))
    end
  end
end
