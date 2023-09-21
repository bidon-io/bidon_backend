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

ActiveRecord::Schema[7.0].define(version: 2023_09_21_073317) do
  # These are extensions that must be enabled in order to support this database
  enable_extension "pgcrypto"
  enable_extension "plpgsql"

  create_table "app_demand_profiles", force: :cascade do |t|
    t.bigint "app_id", null: false
    t.string "account_type", null: false
    t.bigint "account_id", null: false
    t.bigint "demand_source_id", null: false
    t.jsonb "data", default: {}
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.bigint "public_uid"
    t.index ["account_type", "account_id"], name: "index_app_demand_profiles_on_account"
    t.index ["app_id", "demand_source_id"], name: "index_app_demand_profiles_on_app_id_and_demand_source_id", unique: true
    t.index ["app_id"], name: "index_app_demand_profiles_on_app_id"
    t.index ["demand_source_id"], name: "index_app_demand_profiles_on_demand_source_id"
    t.index ["public_uid"], name: "index_app_demand_profiles_on_public_uid", unique: true
  end

  create_table "app_mmp_profiles", force: :cascade do |t|
    t.bigint "app_id", null: false
    t.date "start_date", null: false
    t.integer "mmp_platform", default: 0
    t.bigint "primary_mmp_account"
    t.bigint "secondary_mmp_account"
    t.boolean "get_spend_from_secondary_mmp_account", default: false
    t.integer "primary_mmp_raw_data_source"
    t.integer "secondary_mmp_raw_data_source"
    t.string "adjust_app_token"
    t.string "adjust_s2s_token"
    t.string "adjust_environment"
    t.string "appsflyer_dev_key"
    t.string "appsflyer_app_id"
    t.string "appsflyer_conversion_keys"
    t.string "firebase_config_keys"
    t.integer "firebase_expiration_duration"
    t.boolean "firebase_tracking", default: false
    t.boolean "facebook_tracking", default: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["app_id"], name: "index_app_mmp_profiles_on_app_id"
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
    t.bigint "public_uid"
    t.index ["app_key"], name: "index_apps_on_app_key", unique: true
    t.index ["package_name", "platform_id"], name: "index_apps_on_package_name_and_platform_id", unique: true
    t.index ["public_uid"], name: "index_apps_on_public_uid", unique: true
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
    t.bigint "segment_id"
    t.boolean "external_win_notifications", default: false, null: false
    t.bigint "public_uid"
    t.index ["app_id"], name: "index_auction_configurations_on_app_id"
    t.index ["public_uid"], name: "index_auction_configurations_on_public_uid", unique: true
    t.index ["segment_id"], name: "index_auction_configurations_on_segment_id"
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
    t.string "label"
    t.bigint "public_uid"
    t.index ["demand_source_id"], name: "index_demand_source_accounts_on_demand_source_id"
    t.index ["public_uid"], name: "index_demand_source_accounts_on_public_uid", unique: true
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
    t.integer "width", default: 0, null: false
    t.integer "height", default: 0, null: false
    t.string "format"
    t.bigint "public_uid"
    t.boolean "bidding"
    t.index ["account_type", "account_id"], name: "index_line_items_on_account"
    t.index ["app_id"], name: "index_line_items_on_app_id"
    t.index ["public_uid"], name: "index_line_items_on_public_uid", unique: true
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

  create_table "segments", force: :cascade do |t|
    t.string "name", null: false
    t.text "description", null: false
    t.jsonb "filters", default: [], null: false
    t.boolean "enabled", default: true, null: false
    t.bigint "app_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.integer "priority", default: 0, null: false
    t.bigint "public_uid"
    t.index ["app_id"], name: "index_segments_on_app_id"
    t.index ["public_uid"], name: "index_segments_on_public_uid", unique: true
  end

  create_table "users", force: :cascade do |t|
    t.string "email", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.boolean "is_admin", default: false, null: false
    t.string "password_hash", null: false
    t.bigint "public_uid"
    t.index ["email"], name: "index_users_on_email", unique: true
    t.index ["public_uid"], name: "index_users_on_public_uid", unique: true
  end

  add_foreign_key "app_demand_profiles", "apps"
  add_foreign_key "app_demand_profiles", "demand_source_accounts", column: "account_id"
  add_foreign_key "app_demand_profiles", "demand_sources"
  add_foreign_key "app_mmp_profiles", "apps"
  add_foreign_key "apps", "users"
  add_foreign_key "auction_configurations", "apps"
  add_foreign_key "auction_configurations", "segments"
  add_foreign_key "demand_source_accounts", "demand_sources"
  add_foreign_key "line_items", "apps"
  add_foreign_key "line_items", "demand_source_accounts", column: "account_id"
  add_foreign_key "mmp_accounts", "users"
  add_foreign_key "segments", "apps"
end
