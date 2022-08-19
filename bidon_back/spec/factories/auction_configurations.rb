FactoryBot.define do
  factory :auction_configuration do
    name { 'MyString' }
    app { nil }
    ad_type { 1 }
    rounds { '' }
    status { 1 }
    settings { '' }
    pricefloor { 1 }
  end
end
