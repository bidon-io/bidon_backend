module KafkaLogger
  extend self

  def log(message)
    produce_message(message)
  end

  def log_many(messages)
    messages.each { produce_message(_1) }
  end

  private

  def produce_message(message)
    KafkaProducer.produce(message.to_json, topic: message.topic)
  end
end
