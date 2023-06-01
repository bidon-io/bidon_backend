FactoryBot.define do
  factory :segment do
    name { 'MyString' }
    description { 'MyText' }
    filters { '' }
    enabled { false }
    app_id { '' }
  end
end
