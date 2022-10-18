# frozen_string_literal: true

# Example consumer that prints messages payloads
class CountriesConsumer < ApplicationConsumer
  def model
    Country
  end
end
