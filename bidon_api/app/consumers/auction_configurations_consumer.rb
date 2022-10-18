# frozen_string_literal: true

# Example consumer that prints messages payloads
class AuctionConfigurationsConsumer < ApplicationConsumer
  def model
    AuctionConfiguration
  end
end
