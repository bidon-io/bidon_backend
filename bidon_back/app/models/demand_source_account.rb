class DemandSourceAccount < ApplicationRecord
  belongs_to :demand_source
  belongs_to :user, optional: true

  has_many :line_items, dependent: :restrict_with_exception, foreign_key: :account_id, inverse_of: :account

  def extra=(value)
    if value.is_a?(Hash)
      super(value)
    else
      super(JSON.parse(value.gsub('=>', ':')))
    end
  end
end
