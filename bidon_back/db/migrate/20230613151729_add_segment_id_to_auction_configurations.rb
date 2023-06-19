class AddSegmentIdToAuctionConfigurations < ActiveRecord::Migration[7.0]
  def change
    add_belongs_to :auction_configurations, :segment, foreign_key: true
  end
end
