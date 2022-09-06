class AddWidthAndHeightToAdUnit < ActiveRecord::Migration[7.0]
  def change
    add_column :line_items, :width, :integer, null: false, default: 0
    add_column :line_items, :height, :integer, null: false, default: 0
  end
end
