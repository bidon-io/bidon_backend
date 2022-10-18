# frozen_string_literal: true

# Application consumer from which all Karafka consumers should inherit
# You can rename it if it would conflict with your current code base (in case you're integrating
# Karafka with other frameworks)
class ApplicationConsumer < Karafka::BaseConsumer
  def consume
    messages.each do |message|
      data = message.payload['payload']['after']
      app = model.find_or_initialize_by(id: data['id'])
      app.assign_attributes(params_for(data))
      app.save!
    end
  end

  def params_for(data)
    data.except('id')
  end

  def model
    raise NotImplementedError
  end
end
