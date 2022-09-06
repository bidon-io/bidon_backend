class AuctionConfiguration < Sequel::Model
  plugin :enum

  enum :ad_type, AdType::ENUM
end
