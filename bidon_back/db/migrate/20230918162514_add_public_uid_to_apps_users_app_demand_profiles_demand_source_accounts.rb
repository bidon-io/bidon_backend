class AddPublicUidToAppsUsersAppDemandProfilesDemandSourceAccounts < ActiveRecord::Migration[7.0]
  def change
    add_column :apps, :public_uid, :bigint
    add_index :apps, :public_uid, unique: true

    add_column :users, :public_uid, :bigint
    add_index :users, :public_uid, unique: true

    add_column :app_demand_profiles, :public_uid, :bigint
    add_index :app_demand_profiles, :public_uid, unique: true

    add_column :demand_source_accounts, :public_uid, :bigint
    add_index :demand_source_accounts, :public_uid, unique: true
  end
end
