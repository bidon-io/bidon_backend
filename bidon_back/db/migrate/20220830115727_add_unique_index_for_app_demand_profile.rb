class AddUniqueIndexForAppDemandProfile < ActiveRecord::Migration[7.0]
  def change
    add_index :app_demand_profiles, %i[app_id demand_source_id], unique: true
  end
end
