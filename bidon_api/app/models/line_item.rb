class LineItem < Sequel::Model
  plugin :enum

  many_to_one :demand_source_account, key: :account_id

  enum :ad_type, AdType::ENUM
end
