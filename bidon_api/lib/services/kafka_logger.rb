module KafkaLogger
  module_function

  def log_click(event)
    KafkaProducer.produce(prepare_event(event), topic: ENV.fetch('KAFKA_CLICK_TOPIC'))
  end

  def log_reward(event)
    KafkaProducer.produce(prepare_event(event), topic: ENV.fetch('KAFKA_REWARD_TOPIC'))
  end

  def log_show(event)
    KafkaProducer.produce(prepare_event(event), topic: ENV.fetch('KAFKA_SHOW_TOPIC'))
  end

  def log_stats(event)
    KafkaProducer.produce(prepare_event(event), topic: ENV.fetch('KAFKA_STATS_TOPIC'))
  end

  def prepare_event(event)
    JSON.dump(Utils.smash_hash(event))
  end
end
