package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaConfig struct {
	ClientOpts []kgo.Opt
	Topics     map[event.Topic]string
}

func Kafka() (conf KafkaConfig, err error) {
	seeds := strings.Split(os.Getenv("KAFKA_BROKERS_LIST"), ", ")
	clientID := os.Getenv("KAFKA_CLIENT_ID")

	conf.ClientOpts = []kgo.Opt{
		kgo.AllowAutoTopicCreation(),
		kgo.SeedBrokers(seeds...),
		kgo.ClientID(clientID),
	}

	batchMaxBytes := os.Getenv("KAFKA_BATCH_MAX_BYTES")
	if batchMaxBytes != "" {
		value, err := strconv.Atoi(batchMaxBytes)
		if err != nil {
			return conf, fmt.Errorf("invalid KAFKA_BATCH_MAX_BYTES: %v", err)
		}

		conf.ClientOpts = append(conf.ClientOpts, kgo.ProducerBatchMaxBytes(int32(value)))
	}

	linger := os.Getenv("KAFKA_DELIVERY_INTERVAL")
	if linger != "" {
		value, err := strconv.Atoi(linger)
		if err != nil {
			return conf, fmt.Errorf("invalid KAFKA_DELIVERY_INTERVAL: %v", err)
		}

		conf.ClientOpts = append(conf.ClientOpts, kgo.ProducerLinger(time.Second*time.Duration(value)))
	}

	configTopic := os.Getenv("KAFKA_CONFIG_TOPIC")
	if configTopic == "" {
		return conf, fmt.Errorf("empty KAFKA_CONFIG_TOPIC: %v", err)
	}

	conf.Topics = map[event.Topic]string{
		event.ConfigTopic: configTopic,
	}

	return
}
