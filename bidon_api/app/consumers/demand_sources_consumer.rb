# frozen_string_literal: true

# Example consumer that prints messages payloads
class DemandSourcesConsumer < ApplicationConsumer
  def model
    DemandSource
  end
end
