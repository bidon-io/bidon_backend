class AddBiddingToLineItem < ActiveRecord::Migration[7.0]
  def change
    add_column :line_items, :bidding, :boolean
  end
end
