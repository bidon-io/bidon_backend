class CreateLineItems < ActiveRecord::Migration[7.0]
  def change
    create_table :line_items do |t|
      t.belongs_to :app, null: false, foreign_key: true
      t.belongs_to :account, null: false, polymorphic: true
      t.string :human_name, null: false
      t.string :code, null: false
      t.decimal :bid_floor
      t.integer :ad_type, null: false
      t.jsonb :extra, default: {}

      t.timestamps

      t.foreign_key(:demand_source_accounts, column: :account_id)
    end
  end
end
