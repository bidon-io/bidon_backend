class AddLabelToDemandSoruceAccount < ActiveRecord::Migration[7.0]
  def change
    add_column :demand_source_accounts, :label, :string
  end
end
