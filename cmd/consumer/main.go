package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	model "github.com/indralaksmana/queue_demo/model"
	kafka "github.com/indralaksmana/queue_demo/pkg/kafka"

	"github.com/sirupsen/logrus"
)

func main() {
	// Setup Logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	// Initialize consumer
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	kafkaBroker := []string{os.Getenv("KAFKA_BROKER")}
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaGroup := os.Getenv("KAFKA_GROUP")

	/**
	 * Setup a new queue consumer group
	 */
	consumer := Consumer{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := kafka.ConnectConsumer(kafkaBroker, kafkaGroup)
	if err != nil {
		logrus.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, []string{kafkaTopic}, &consumer); err != nil {
				logrus.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	logrus.Println("Queue consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		logrus.Println("terminating: context cancelled")
	case <-sigterm:
		logrus.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		logrus.Panicf("Error closing client: %v", err)
	}

}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		var textMessage model.Message
		json.Unmarshal([]byte(message.Value), &textMessage)

		logrus.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", textMessage, message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}

	return nil
}
