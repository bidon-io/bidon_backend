# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[7.0].define(version: 2022_08_18_160109) do
  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"

  create_table "app_demand_profiles", force: :cascade do |t|
    t.bigint "app_id", null: false
    t.string "account_type", null: false
    t.bigint "account_id", null: false
    t.bigint "demand_source_id", null: false
    t.jsonb "data", default: {}
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["account_type", "account_id"], name: "index_app_demand_profiles_on_account"
    t.index ["app_id"], name: "index_app_demand_profiles_on_app_id"
    t.index ["demand_source_id"], name: "index_app_demand_profiles_on_demand_source_id"
  end

  create_table "apps", force: :cascade do |t|
    t.bigint "user_id", null: false
    t.integer "platform_id", null: false
    t.string "human_name", null: false
    t.string "package_name"
    t.string "app_key"
    t.jsonb "settings", default: {}
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["app_key"], name: "index_apps_on_app_key", unique: true
    t.index ["package_name"], name: "index_apps_on_package_name", unique: true
    t.index ["user_id"], name: "index_apps_on_user_id"
  end

  create_table "auction_configurations", force: :cascade do |t|
    t.string "name"
    t.bigint "app_id", null: false
    t.integer "ad_type", null: false
    t.jsonb "rounds", default: []
    t.integer "status"
    t.jsonb "settings", default: {}
    t.float "pricefloor", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["app_id"], name: "index_auction_configurations_on_app_id"
  end

  create_table "countries", force: :cascade do |t|
    t.string "alpha_2_code", null: false
    t.string "alpha_3_code", null: false
    t.string "human_name"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["alpha_2_code"], name: "index_countries_on_alpha_2_code", unique: true
    t.index ["alpha_3_code"], name: "index_countries_on_alpha_3_code", unique: true
  end

  create_table "demand_source_accounts", force: :cascade do |t|
    t.bigint "demand_source_id", null: false
    t.bigint "user_id", null: false
    t.string "type", null: false
    t.jsonb "extra", default: {}
    t.boolean "bidding", default: false
    t.boolean "is_default"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["demand_source_id"], name: "index_demand_source_accounts_on_demand_source_id"
  end

  create_table "demand_sources", force: :cascade do |t|
    t.string "api_key", null: false
    t.string "human_name", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["api_key"], name: "index_demand_sources_on_api_key", unique: true
  end

  create_table "line_items", force: :cascade do |t|
    t.bigint "app_id", null: false
    t.string "account_type", null: false
    t.bigint "account_id", null: false
    t.string "human_name", null: false
    t.string "code", null: false
    t.decimal "bid_floor"
    t.integer "ad_type", null: false
    t.jsonb "extra", default: {}
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["account_type", "account_id"], name: "index_line_items_on_account"
    t.index ["app_id"], name: "index_line_items_on_app_id"
  end

  create_table "mmp_accounts", force: :cascade do |t|
    t.bigint "user_id", null: false
    t.string "human_name", null: false
    t.integer "account_type", null: false
    t.boolean "use_s3", default: false
    t.string "s3_access_key_id"
    t.string "s3_secret_access_key"
    t.string "s3_bucket_name"
    t.string "s3_region"
    t.string "s3_home_folder"
    t.string "master_api_token"
    t.string "user_token"
    t.boolean "is_global_account", default: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["user_id"], name: "index_mmp_accounts_on_user_id"
  end

  create_table "users", force: :cascade do |t|
    t.string "email", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["email"], name: "index_users_on_email", unique: true
  end

  add_foreign_key "app_demand_profiles", "apps"
  add_foreign_key "app_demand_profiles", "demand_source_accounts", column: "account_id"
  add_foreign_key "app_demand_profiles", "demand_sources"
  add_foreign_key "apps", "users"
  add_foreign_key "auction_configurations", "apps"
  add_foreign_key "demand_source_accounts", "demand_sources"
  add_foreign_key "line_items", "apps"
  add_foreign_key "line_items", "demand_source_accounts", column: "account_id"
  add_foreign_key "mmp_accounts", "users"
end
