class CreateAuctionConfigurations < ActiveRecord::Migration[7.0]
  def change
    create_table :auction_configurations do |t|
      t.string :name
      t.references :app, null: false, foreign_key: true
      t.integer :ad_type, null: false
      t.jsonb :rounds, default: []
      t.integer :status
      t.jsonb :settings, default: {}
      t.float :pricefloor, null: false

      t.timestamps
    end
  end
end
