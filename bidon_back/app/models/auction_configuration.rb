class AuctionConfiguration < ApplicationRecord
  belongs_to :app

  enum ad_type: AdType::ENUM

  validates :name, :pricefloor, :ad_type, :rounds, presence: true
end
