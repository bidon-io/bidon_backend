class AuctionConfiguration < ApplicationRecord
  belongs_to :app

  enum ad_type: AdType::ENUM

  def rounds=(value)
    if value.is_a?(Array)
      super(value)
    else
      super(JSON.parse(value.gsub('=>', ':')))
    end
  end
end
