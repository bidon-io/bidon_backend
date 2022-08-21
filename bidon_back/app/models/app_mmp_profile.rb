class AppMmpProfile < ApplicationRecord
  belongs_to :app

  belongs_to :mmp_account_primary,
             class_name:  'MmpAccount',
             foreign_key: 'primary_mmp_account',
             inverse_of:  :primary_app_profiles

  belongs_to :mmp_account_secondary,
             class_name:  'MmpAccount',
             foreign_key: 'secondary_mmp_account',
             inverse_of:  :secondary_app_profiles,
             optional:    true

  enum mmp_platform: { none: 0, appsflyer: 1, adjust: 2 }, _prefix: 'mmp_platform'
  enum primary_mmp_raw_data_source:   { none: 0, appsflyer: 1, adjust: 2 }, _prefix: 'primary_raw_source'
  enum secondary_mmp_raw_data_source: { none: 0, appsflyer: 1, adjust: 2 }, _prefix: 'secondary_raw_source'

  validates :start_date, presence: true

  # t.string :adjust_app_token
  # t.string :adjust_s2s_token
  # t.string :adjust_environment

  # t.string :appsflyer_dev_key
  # t.string :appsflyer_app_id
  # t.string :appsflyer_conversion_keys

  # t.string :firebase_config_keys
  # t.integer :firebase_expiration_duration
  # t.boolean :firebase_tracking, default: false

  # t.boolean :facebook_tracking, default: false
end
