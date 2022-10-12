class CreateDemandSourceAccounts < ActiveRecord::Migration[7.0]
  def change
    create_table :demand_source_accounts do |t|
      t.belongs_to :demand_source, null: false, foreign_key: true
      t.bigint :user_id, null: false
      t.string :type, null: false
      t.jsonb :extra, default: {}
      t.boolean :bidding, default: false
      t.boolean :is_default

      t.timestamps
    end
  end
end
