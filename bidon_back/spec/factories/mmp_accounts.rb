FactoryBot.define do
  factory :mmp_account do
    user { nil }
    human_name { 'MyString' }
    account_type { 1 }
    use_s3 { false }
    s3_access_key_id { 'MyString' }
    s3_secret_access_key { 'MyString' }
    s3_bucket_name { 'MyString' }
    s3_region { 'MyString' }
    s3_home_folder { 'MyString' }
    master_api_token { 'MyString' }
    user_token { 'MyString' }
    is_global_account { false }
  end
end
