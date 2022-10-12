class AuctionConfiguration < ApplicationRecord
  enum ad_type: AdType::ENUM
end
