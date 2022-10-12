class CreateCountries < ActiveRecord::Migration[7.0]
  def change
    create_table :countries do |t|
      t.string :alpha_2_code, null: false, index: { unique: true }
      t.string :alpha_3_code, null: false, index: { unique: true }
      t.string :human_name

      t.timestamps
    end
  end
end
