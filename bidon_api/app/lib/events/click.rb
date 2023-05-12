# frozen_string_literal: true

class Events::Click < Events::Base
  KAFKA_TOPIC = ENV.fetch('KAFKA_CLICK_TOPIC')
end
