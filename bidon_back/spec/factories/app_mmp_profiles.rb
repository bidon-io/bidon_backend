FactoryBot.define do
  factory :app_mmp_profile do
    app { nil }
    start_date { '2022-08-18' }
    mmp_platform { 1 }
    primary_mmp_account { '' }
    secondary_mmp_account { '' }
    get_spend_from_secondary_mmp_account { false }
    primary_mmp_raw_data_source { 1 }
    secondary_mmp_raw_data_source { 1 }
    adjust_app_token { 'MyString' }
    adjust_s2s_token { 'MyString' }
    adjust_environment { 'MyString' }
    appsflyer_dev_key { 'MyString' }
    appsflyer_app_id { 'MyString' }
    appsflyer_conversion_keys { 'MyString' }
    firebase_config_keys { 'MyString' }
    firebase_expiration_duration { 1 }
    firebase_tracking { false }
    facebook_tracking { false }
  end
end
