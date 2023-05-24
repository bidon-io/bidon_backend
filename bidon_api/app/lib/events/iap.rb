# frozen_string_literal: true

class Events::Iap < Events::Base
  KAFKA_TOPIC = ENV.fetch('KAFKA_IAP_TOPIC')
end
