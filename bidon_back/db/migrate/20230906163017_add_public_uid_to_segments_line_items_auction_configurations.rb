class AddPublicUidToSegmentsLineItemsAuctionConfigurations < ActiveRecord::Migration[7.0]
  def change
    add_column :segments, :public_uid, :bigint
    add_index :segments, :public_uid, unique: true

    add_column :line_items, :public_uid, :bigint
    add_index :line_items, :public_uid, unique: true

    add_column :auction_configurations, :public_uid, :bigint
    add_index :auction_configurations, :public_uid, unique: true
  end
end
