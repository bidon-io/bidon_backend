class User < ApplicationRecord
  has_many :apps, dependent: :restrict_with_exception
  has_many :demand_source_accounts, dependent: :restrict_with_exception
  has_many :mmp_accounts, dependent: :restrict_with_exception

  validates :email, presence: true, uniqueness: true
end
