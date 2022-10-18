# frozen_string_literal: true

# Example consumer that prints messages payloads
class AppDemandProfilesConsumer < ApplicationConsumer
  def model
    AppDemandProfile
  end
end
