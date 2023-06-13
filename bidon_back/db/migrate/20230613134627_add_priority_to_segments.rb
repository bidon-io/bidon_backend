class AddPriorityToSegments < ActiveRecord::Migration[7.0]
  def change
    add_column :segments, :priority, :integer, null: false, default: 0
  end
end
