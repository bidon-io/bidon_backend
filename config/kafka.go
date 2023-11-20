package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Topic string

const (
	AdEventsTopic Topic = "ad_events"
)

type KafkaConfig struct {
	ClientOpts []kgo.Opt
	Topics     map[Topic]string
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

	conf.Topics = map[Topic]string{
		AdEventsTopic: os.Getenv("KAFKA_AD_EVENTS_TOPIC"),
	}

	return
}
