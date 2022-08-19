class AuctionConfiguration < ApplicationRecord
  belongs_to :app

  enum ad_type: { interstitial: 1, banner: 2, rewarded_video: 3 }

  validates :name, :pricefloor, :ad_type, :rounds, presence: true
end
