class AddExternalNotificationToAuctionConfigurations < ActiveRecord::Migration[7.0]
  def change
    add_column :auction_configurations, :external_win_notifications, :boolean, default: false, null: false
  end
end
