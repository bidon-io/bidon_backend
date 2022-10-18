# frozen_string_literal: true

# Example consumer that prints messages payloads
class UsersConsumer < ApplicationConsumer
  def model
    User
  end
end
