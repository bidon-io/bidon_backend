FactoryBot.define do
  factory :demand_source_account do
    demand_source { nil }
    user { nil }
    type { '' }
    extra { '' }
    bidding { false }
  end
end
