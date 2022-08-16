FactoryBot.define do
  factory :app do
    user { nil }
    platform_id { 1 }
    human_name { 'MyString' }
    package_name { 'MyString' }
    app_key { 'MyString' }
    settings { '' }
  end
end
