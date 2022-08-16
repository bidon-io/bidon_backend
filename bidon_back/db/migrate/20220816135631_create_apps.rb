class CreateApps < ActiveRecord::Migration[7.0]
  def change
    create_table :apps do |t|
      t.belongs_to :user, null: false, foreign_key: true
      t.integer :platform_id, null: false
      t.string :human_name, null: false
      t.string :package_name, index: { unique: true }
      t.string :app_key, index: { unique: true }
      t.jsonb :settings, default: {}

      t.timestamps
    end
  end
end
