FactoryBot.define do
  factory :line_item do
    app { nil }
    account { nil }
    human_name { 'MyString' }
    code { 'MyString' }
    bid_floor { '9.99' }
    ad_type { 1 }
    extra { '' }
  end
end
