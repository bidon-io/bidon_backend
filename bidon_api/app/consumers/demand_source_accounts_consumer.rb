# frozen_string_literal: true

# Example consumer that prints messages payloads
class DemandSourceAccountsConsumer < ApplicationConsumer
  def model
    DemandSourceAccount
  end
end
