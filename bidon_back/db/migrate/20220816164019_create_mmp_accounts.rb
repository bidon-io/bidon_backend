class CreateMmpAccounts < ActiveRecord::Migration[7.0]
  def change
    create_table :mmp_accounts do |t|
      t.belongs_to :user, null: false, foreign_key: true
      t.string :human_name, null: false
      t.integer :account_type, null: false
      t.boolean :use_s3, default: false
      t.string :s3_access_key_id
      t.string :s3_secret_access_key
      t.string :s3_bucket_name
      t.string :s3_region
      t.string :s3_home_folder
      t.string :master_api_token
      t.string :user_token
      t.boolean :is_global_account, default: false

      t.timestamps
    end
  end
end
