class User < ApplicationRecord
  has_many :apps, dependent: :restrict_with_exception

  validates :email, presence: true, uniqueness: true
end
