package kafka

import (
	"time"

	"github.com/Shopify/sarama"
)

func ConnectProducer(kafkaBrokerUrls []string) (sarama.SyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	conn, err := sarama.NewSyncProducer(kafkaBrokerUrls, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ConnectConsumer(kafkaBrokerUrls []string, kafkaGroup string) (sarama.ConsumerGroup, error) {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Retry.Backoff = 3
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	// Create new consumer
	conn, err := sarama.NewConsumerGroup(kafkaBrokerUrls, kafkaGroup, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
