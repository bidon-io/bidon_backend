# frozen_string_literal: true

# Example consumer that prints messages payloads
class LineItemsConsumer < ApplicationConsumer
  def model
    LineItem
  end
end
