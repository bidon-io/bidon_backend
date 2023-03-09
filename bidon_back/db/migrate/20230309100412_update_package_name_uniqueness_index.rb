class UpdatePackageNameUniquenessIndex < ActiveRecord::Migration[7.0]
  def up
    remove_index :apps, :package_name, name: 'index_apps_on_package_name'
    add_index :apps, %i[package_name platform_id], unique: true, name: 'index_apps_on_package_name_and_platform_id'
  end

  def down
    remove_index :apps, %i[package_name platform_id], name: 'index_apps_on_package_name_and_platform_id'
    add_index :apps, :package_name, unique: true, name: 'index_apps_on_package_name'
  end
end
