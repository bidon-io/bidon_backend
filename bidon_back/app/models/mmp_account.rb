class MmpAccount < ApplicationRecord
  belongs_to :user

  has_many :primary_app_profiles,
           class_name:  'AppMmpProfile',
           foreign_key: 'primary_mmp_account',
           inverse_of:  :mmp_account_primary,
           dependent:   :restrict_with_exception

  has_many :secondary_app_profiles,
           class_name:  'AppMmpProfile',
           foreign_key: 'secondary_mmp_account',
           inverse_of:  :mmp_account_secondary,
           dependent:   :restrict_with_exception

  enum account_type: { appsflyer: 1, adjust: 2 }

  validates :human_name, :account_type, presence: true
end
