package broker

import (
	"encoding/json"
	"log"
	"os"
	"price-input-storage/model"
	"price-input-storage/repository"
	"time"

	"github.com/IBM/sarama"
)

type (
	MessageBroker interface {
		Consume(topic string) error
	}

	messageBrokerImpl struct {
		repo repository.PriceInputRepo
	}
)

func NewMessageBroker(repo repository.PriceInputRepo) MessageBroker {
	return &messageBrokerImpl{repo: repo}
}

func (m *messageBrokerImpl) Consume(topic string) (err error) {
	config := sarama.NewConfig()

	consumer, err := sarama.NewConsumer([]string{os.Getenv("KAFKA_HOST")}, config)
	if err != nil {
		log.Fatal("NewConsumer err: ", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("ConsumePartition err: ", err)
	}
	defer partitionConsumer.Close()
	for message := range partitionConsumer.Messages() {
		var price model.Price
		err = json.Unmarshal(message.Value, &price)
		if err != nil {
			log.Fatal("Unmarshalling err: ", err)
		}
		log.Printf("[price-input-storeage-consumer] partitionid: %d; offset:%d, value: %v\n", message.Partition, message.Offset, price)

		price.Id = string(message.Key)
		price.Date = time.Now().Unix()

		m.repo.Insert(price)
	}

	return err
}
