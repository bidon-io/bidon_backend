class CreateDemandSources < ActiveRecord::Migration[7.0]
  def change
    create_table :demand_sources do |t|
      t.string :api_key, null: false, index: { unique: true }
      t.string :human_name, null: false

      t.timestamps
    end
  end
end
