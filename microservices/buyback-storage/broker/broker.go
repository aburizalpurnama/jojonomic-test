package broker

import (
	"buyback-service/usecases/request"
	"encoding/json"
	"log"
	"os"
	"topup-storage/model"
	"topup-storage/repository"

	"github.com/IBM/sarama"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

type (
	MessageBroker interface {
		Consume(topic string) error
	}

	messageBrokerImpl struct {
		repo repository.TopupRepo
	}
)

func NewMessageBroker(repo repository.TopupRepo) MessageBroker {
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
		var data request.BuybackPayload
		err = json.Unmarshal(message.Value, &data)
		if err != nil {
			logrus.Error("Unmarshalling err: ", err)
		}
		logrus.Infof("[topup-storeage-consumer] partitionid: %d; offset:%d, value: %v\n", message.Partition, message.Offset, data)

		var (
			trx model.Transaction
			acc model.Account
		)

		err = copier.Copy(&trx, &data.Transaction)
		if err != nil {
			return err
		}

		err = copier.Copy(&acc, &data.Account)
		if err != nil {
			return err
		}

		err = m.repo.InsertTransaksiAndUpdateBalance(trx, acc)
		if err != nil {
			return err
		}
	}

	return err
}
