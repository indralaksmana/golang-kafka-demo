package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func PushToQueue(kafkaBrokerUrls []string, clientId string, topic string, message []byte) error {

	producer, err := ConnectProducer(kafkaBrokerUrls)
	if err != nil {
		return err
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}
