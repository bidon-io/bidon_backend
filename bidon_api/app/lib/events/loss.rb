# frozen_string_literal: true

class Events::Loss < Events::Base
  KAFKA_TOPIC = ENV.fetch('KAFKA_LOSS_TOPIC')
end
