class CreateAppMmpProfiles < ActiveRecord::Migration[7.0]
  def change # rubocop:disable Metrics/AbcSize, Metrics/MethodLength
    create_table :app_mmp_profiles do |t|
      t.belongs_to :app, null: false, foreign_key: true
      t.date :start_date, null: false
      t.integer :mmp_platform, default: 0
      t.bigint :primary_mmp_account
      t.bigint :secondary_mmp_account
      t.boolean :get_spend_from_secondary_mmp_account, default: false
      t.integer :primary_mmp_raw_data_source
      t.integer :secondary_mmp_raw_data_source
      t.string :adjust_app_token
      t.string :adjust_s2s_token
      t.string :adjust_environment
      t.string :appsflyer_dev_key
      t.string :appsflyer_app_id
      t.string :appsflyer_conversion_keys
      t.string :firebase_config_keys
      t.integer :firebase_expiration_duration
      t.boolean :firebase_tracking, default: false
      t.boolean :facebook_tracking, default: false

      t.timestamps
    end
  end
end
