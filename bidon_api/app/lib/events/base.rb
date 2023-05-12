# frozen_string_literal: true

# Events::Base is a base class for all events.
# This is the final stop before events are serialized and sent to Kafka.
# It provides a way for subclasses to specify the Kafka topic to send the event to, and optionally change the payload
# before it is sent to kafka.
class Events::Base
  # @param [EventParams] params
  def initialize(params)
    @params = params
  end

  def topic
    self.class::KAFKA_TOPIC
  end

  # ActiveSupport uses this method to serialize the object
  # with `as_json` and then serialize the result with `to_json`
  def to_hash
    hash = @params.to_hash

    hash['show'] = hash['bid'] unless hash.key?('show')

    Utils.smash_hash(hash)
  end
end
