# frozen_string_literal: true

class Events::Reward < Events::Base
  KAFKA_TOPIC = ENV.fetch('KAFKA_REWARD_TOPIC')
end
