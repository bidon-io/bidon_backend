class AddPasswordAndIsAdminToUsers < ActiveRecord::Migration[7.0]
  def change
    add_column :users, :is_admin, :boolean, null: false, default: false

    add_column :users, :password_hash, :string
    User.update_all(password_hash: 'password') # rubocop:disable Rails/SkipsModelValidations
    change_column_null :users, :password_hash, false
  end
end
