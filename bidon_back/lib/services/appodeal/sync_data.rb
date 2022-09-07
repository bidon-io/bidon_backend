module Appodeal
  class SyncData
    APP_IDS = [
      454_236, # iOS test app
      726_191, # Android test app
    ].freeze

    DEMANDS = [
      20,   # AdMob
      50,   # Applovin
      1260, # BidMachine
    ].freeze

    APPLOVIN_DEMAND = 50
    APPLOVIN_ACCOUNT_ID = 4

    attr_reader :app_ids, :demand_ids, :appodeal_connection

    # @param [Array<Integer>] app_ids
    # @param [Array<Integer>] demand_ids
    # @param [Class<ActiveRecord::Base>] appodeal_connection
    def initialize(app_ids: APP_IDS, demand_ids: DEMANDS, appodeal_connection: AppodealPg)
      @app_ids = app_ids
      @demand_ids = demand_ids
      @appodeal_connection = appodeal_connection
    end

    def call
      ApplicationRecord.transaction do
        sync_users
        sync_apps
        sync_countries

        sync_demand_sources # campaign_types
        sync_demand_source_accounts # network_accounts
        sync_app_demand_profiles # app_network_profiles
        sync_line_items # ad_units

        sync_mmp_accounts
        sync_app_mmp_profiles # app_attributions
      end
    end

    private

    def sync_users
      users = appodeal_connection.execute <<~SQL.squish
        SELECT id, email FROM users WHERE id IN (
          SELECT user_id FROM apps WHERE #{app_ids_filter}
        )
      SQL

      build_models(User, users)
    end

    def sync_apps
      apps = appodeal_connection.execute <<~SQL.squish
        SELECT id, user_id, platform_id, name AS human_name, package_name, app_key FROM apps WHERE #{app_ids_filter}
      SQL

      build_models(App, apps)
    end

    def sync_countries
      countries = appodeal_connection.execute <<~SQL.squish
        SELECT id, name AS human_name, code AS alpha_2_code, alpha_3_code FROM countries
      SQL

      build_models(Country, countries)
    end

    def sync_demand_sources
      demand_sources = appodeal_connection.execute <<~SQL.squish
        SELECT id, human_name, api_name AS api_key FROM campaign_types WHERE id IN (#{demand_ids.join(', ')})
      SQL

      build_models(DemandSource, demand_sources)
    end

    def sync_demand_source_accounts
      accounts = appodeal_connection.execute <<~SQL.squish
        SELECT id, extra, bidding_account AS bidding, owner_id AS user_id,
        CASE
          WHEN type = 'BidmachineAccount' THEN 'DemandSourceAccount::BidMachine'
          WHEN type = 'AdmobAccount' THEN 'DemandSourceAccount::Admob'
          WHEN type = 'ApplovinAccount' THEN 'DemandSourceAccount::Applovin'
        ELSE NULL END AS type,
        CASE
          WHEN type = 'BidmachineAccount' THEN (SELECT id from campaign_types WHERE api_name = 'appodeal_exchange')
          WHEN type = 'AdmobAccount' THEN (SELECT id from campaign_types WHERE api_name = 'admob')
          WHEN type = 'ApplovinAccount' THEN (SELECT id from campaign_types WHERE api_name = 'applovin')
        ELSE NULL END AS demand_source_id
        FROM network_accounts
        WHERE type IN ('AdmobAccount', 'BidmachineAccount', 'ApplovinAccount')
          AND (
            owner_id = 0
            OR owner_id IN (SELECT user_id FROM apps WHERE #{app_ids_filter})
          )
      SQL

      build_models(DemandSourceAccount, accounts)
    end

    def sync_app_demand_profiles
      profiles = appodeal_connection.execute <<~SQL.squish
        SELECT DISTINCT ON (app_id, demand_source_id)
               id, app_id, network AS demand_source_id, data::jsonb, account_id,
          CASE
            WHEN network = 1260 THEN 'DemandSourceAccount::BidMachine'
            WHEN network = 20 THEN 'DemandSourceAccount::Admob'
            WHEN network = 50 THEN 'DemandSourceAccount::Applovin'
          ELSE NULL END AS account_type
        FROM app_network_profiles
        WHERE app_id IN (#{app_ids.join(', ')})
          AND (
            (network IN (#{(demand_ids - [APPLOVIN_DEMAND]).join(', ')}))
            OR
            (network = #{APPLOVIN_DEMAND} AND account_id = #{APPLOVIN_ACCOUNT_ID})
          )
        ORDER BY app_id, demand_source_id, updated_at DESC
      SQL

      build_models(AppDemandProfile, profiles)
    end

    def sync_line_items
      legacy_ad_types = { interstitial: 0, banner: 1, rewarded: 5 }

      profiles = appodeal_connection.execute <<~SQL.squish
        SELECT id, app_id, bid_floor, account_id, ad_type, code, extra, height, width,
        COALESCE(label, package_name) AS human_name,
        CASE
          WHEN ad_type = #{legacy_ad_types[:interstitial]} THEN #{AdType::ENUM[:interstitial]}
          WHEN ad_type = #{legacy_ad_types[:banner]} THEN #{AdType::ENUM[:banner]}
          WHEN ad_type = #{legacy_ad_types[:rewarded]} THEN #{AdType::ENUM[:rewarded]}
        END AS ad_type,
        CASE
          WHEN account_type = 'BidmachineAccount' THEN 'DemandSourceAccount::BidMachine'
          WHEN account_type = 'AdmobAccount' THEN 'DemandSourceAccount::Admob'
          WHEN account_type = 'ApplovinAccount' THEN 'DemandSourceAccount::Applovin'
        ELSE NULL END AS account_type
        FROM ad_units
        WHERE app_id IN (#{app_ids.join(', ')})
          AND ad_type IN (#{legacy_ad_types.values.join(', ')})
          AND account_type IN ('BidmachineAccount', 'AdmobAccount', 'ApplovinAccount')
      SQL

      build_models(LineItem, profiles)
    end

    def sync_mmp_accounts
      mmp_accounts = appodeal_connection.execute <<~SQL.squish
        SELECT id, #{appodeal_user.id} AS user_id, name AS human_name, account_type, use_s3,
          #{decrypt('s3_access_key_id')} AS s3_access_key_id,
          #{decrypt('s3_secret_access_key')} AS s3_secret_access_key,
          #{decrypt('s3_bucket_name')} AS s3_bucket_name,
          #{decrypt('s3_region')} AS s3_region,
          CASE WHEN s3_home_folder = '' THEN '' ELSE #{decrypt('s3_home_folder')} END AS s3_home_folder,
          CASE WHEN master_api_token = '' THEN '' ELSE #{decrypt('master_api_token')} END AS master_api_token,
          CASE WHEN user_token = '' THEN '' ELSE #{decrypt('user_token')} END AS user_token,
          is_appodeal_account AS is_global_account
        FROM mmp_accounts WHERE is_appodeal_account = true
      SQL

      build_models(MmpAccount, mmp_accounts)
    end

    def sync_app_mmp_profiles # rubocop:disable Metrics/MethodLength
      app_mmp_profiles = appodeal_connection.execute <<~SQL.squish
        SELECT ap.id, ap.app_id,
          coalesce(aa.date, CURRENT_DATE) AS start_date,
          CASE
            WHEN ap.holistic_solution_attribution_platform = 2 THEN 1
            WHEN ap.holistic_solution_attribution_platform = 1 THEN 2
            ELSE 0
          END AS mmp_platform,
          aa.primary_mmp_account,
          aa.secondary_mmp_account,
          aa.get_spend_from_secondary_mmp_account,
          aa.primary_mmp_raw_data_source,
          aa.secondary_mmp_raw_data_source,
          CASE
            WHEN ap.encrypted_holistic_solution_adjust_app_token = '' THEN ''
            ELSE #{decrypt('ap.encrypted_holistic_solution_adjust_app_token')}
          END AS adjust_app_token,
          CASE
            WHEN ap.encrypted_holistic_solution_adjust_s2s_token = '' THEN ''
            ELSE #{decrypt('ap.encrypted_holistic_solution_adjust_s2s_token')}
          END AS adjust_s2s_token,
          ap.holistic_solution_adjust_environment AS adjust_environment,
          CASE
            WHEN ap.encrypted_holistic_solution_appsflyer_dev_key = '' THEN ''
            ELSE #{decrypt('ap.encrypted_holistic_solution_appsflyer_dev_key')}
          END AS appsflyer_dev_key,
          ap.holistic_solution_appsflyer_app_id AS appsflyer_app_id,
          ap.holistic_solution_appsflyer_conversion_keys AS appsflyer_conversion_keys,
          ap.holistic_solution_appsflyer_conversion_keys AS appsflyer_conversion_keys,
          ap.holistic_solution_firebase_config_keys AS firebase_config_keys,
          ap.holistic_solution_firebase_expiration_duration AS firebase_expiration_duration,
          ap.holistic_solution_firebase_tracking AS firebase_tracking,
          ap.holistic_solution_facebook_tracking AS facebook_tracking
        FROM app_profiles ap
        LEFT JOIN (
          SELECT * FROM app_attributions
          WHERE app_id IN (#{app_ids.join(', ')}) AND (app_id, date) IN (
            select app_id, max(date) from app_attributions group by app_id
          )
        ) aa ON ap.app_id = aa.app_id
        WHERE ap.app_id IN (#{app_ids.join(', ')})
      SQL

      build_models(AppMmpProfile, app_mmp_profiles)
    end

    def build_models(klass, hashes)
      hashes.each { |hash| build_model(klass, hash) }

      klass.connection.execute("SELECT setval('#{klass.sequence_name}', (SELECT MAX(id) + 1 FROM #{klass.table_name}))")
    end

    def build_model(klass, hash)
      json_cols = klass.columns.select { |col| col.type == :jsonb }.map(&:name)
      json_cols.each do |col|
        next if hash[col].nil?

        hash[col] = JSON.parse(hash[col])
      end

      klass.find_or_initialize_by(hash.slice('id')).update!(hash.compact.except('id'))
    end

    def app_ids_filter
      @app_ids_filter ||= "id IN (#{app_ids.join(', ')})"
    end

    def decrypt(col)
      "pgp_sym_decrypt(#{col}::bytea, '#{secret_key}')"
    end

    def secret_key
      @secret_key ||= ENV.fetch('APPODEAL_PG_SECRET')
    end

    def appodeal_user
      @appodeal_user ||= User.create_or_find_by(id: 1, email: 'admins@appodeal.com')
    end
  end
end
