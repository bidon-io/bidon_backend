# frozen_string_literal: true

class Events::Show < Events::Base
  KAFKA_TOPIC = ENV.fetch('KAFKA_SHOW_TOPIC')
end
