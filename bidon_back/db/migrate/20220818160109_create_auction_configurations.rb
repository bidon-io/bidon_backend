class CreateAuctionConfigurations < ActiveRecord::Migration[7.0]
  def change
    create_table :auction_configurations do |t|
      t.string :name
      t.references :app, null: false, foreign_key: true
      t.integer :ad_type
      t.jsonb :rounds
      t.integer :status
      t.jsonb :settings

      t.timestamps
    end
  end
end
