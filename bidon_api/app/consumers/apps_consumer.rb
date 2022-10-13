# frozen_string_literal: true

# Example consumer that prints messages payloads
class AppsConsumer < ApplicationConsumer
  def consume
    messages.each do |message|
      data = message.payload['payload']['after']
      app = App.find_or_initialize_by(id: data['id'])
      app.assign_attributes(data.except('id'))
      app.save!
    end
  end
end
