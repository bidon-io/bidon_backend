class DemandSourceAccount < ApplicationRecord
  belongs_to :demand_source
  belongs_to :user

  has_many :line_items, dependent: :restrict_with_exception
end
