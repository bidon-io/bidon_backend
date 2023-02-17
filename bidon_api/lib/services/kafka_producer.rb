class KafkaProducer
  include Singleton
  prepend MemoWise

  class KafkaProducerStub
    def produce(_event, _options)
      # do nothing
    end

    def shutdown
      # do nothing
    end
  end

  class << self
    delegate :produce, to: :instance
  end

  delegate :produce, to: :kafka_producer

  def kafka_producer
    delivery_threshold = ENV.fetch('KAFKA_DELIVERY_THRESHOLD').to_i
    delivery_interval = ENV.fetch('KAFKA_DELIVERY_INTERVAL').to_i

    if Rails.env.production?
      producer = kafka.async_producer(delivery_threshold:, delivery_interval:)
      at_exit { producer.shutdown }
      producer
    else
      KafkaProducerStub.new
    end
  end
  memo_wise :kafka_producer

  def kafka
    seed_brokers = ENV.fetch('KAFKA_BROKERS_LIST').split(', ')
    client_id = ENV.fetch('KAFKA_CLIENT_ID')

    Kafka.new(seed_brokers:, client_id:, logger: Rails.logger)
  end
  memo_wise :kafka
end
