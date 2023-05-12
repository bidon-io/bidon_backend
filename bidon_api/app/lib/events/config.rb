# frozen_string_literal: true

class Events::Config < Events::Base
  KAFKA_TOPIC = ENV.fetch('KAFKA_CONFIG_TOPIC')
end
