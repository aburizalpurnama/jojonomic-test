package broker

import (
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

type (
	MessageBroker interface {
		ProcudeData(topic, key string, value []byte) error
	}

	messageBrokerImpl struct{}
)

func NewMessageBroker() MessageBroker {
	return &messageBrokerImpl{}
}

func (m *messageBrokerImpl) ProcudeData(topic, key string, value []byte) (err error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = 5
	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_HOST")}, config)
	if err != nil {
		return fmt.Errorf("failed to initialize NewSyncProducer, err: %w", err)
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{Topic: topic, Key: sarama.StringEncoder(key), Value: sarama.ByteEncoder(value)}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to produce message, err: %w", err)
	}

	log.Printf("[input-price-producer] partition id: %d; offset:%d\n", partition, offset)
	return
}
