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

	conf.Topics = map[event.Topic]string{
		event.ConfigTopic: os.Getenv("KAFKA_CONFIG_TOPIC"),
		event.ShowTopic:   os.Getenv("KAFKA_SHOW_TOPIC"),
		event.ClickTopic:  os.Getenv("KAFKA_CLICK_TOPIC"),
		event.RewardTopic: os.Getenv("KAFKA_REWARD_TOPIC"),
		event.StatsTopic:  os.Getenv("KAFKA_STATS_TOPIC"),
	}

	return
}
