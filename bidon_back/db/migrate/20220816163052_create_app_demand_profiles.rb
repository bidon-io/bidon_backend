class CreateAppDemandProfiles < ActiveRecord::Migration[7.0]
  def change
    create_table :app_demand_profiles do |t|
      t.belongs_to :app, null: false, foreign_key: true
      t.belongs_to :account, null: false, polymorphic: true
      t.belongs_to :demand_source, null: false, foreign_key: true
      t.jsonb :data, default: {}

      t.timestamps

      t.foreign_key(:demand_source_accounts, column: :account_id)
    end
  end
end
