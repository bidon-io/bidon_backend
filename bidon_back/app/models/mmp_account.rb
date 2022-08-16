class MmpAccount < ApplicationRecord
  belongs_to :user

  enum account_type: { appsflyer: 1, adjust: 2 }

  validates :human_name, :account_type, presence: true
end
