package broker

import (
	"encoding/json"
	"log"
	"os"
	"time"
	svcUsecases "topup-service/usecases"
	"topup-storage/model"
	"topup-storage/repository"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
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
		var topup svcUsecases.Topup
		err = json.Unmarshal(message.Value, &topup)
		if err != nil {
			logrus.Error("Unmarshalling err: ", err)
		}
		logrus.Infof("[topup-storeage-consumer] partitionid: %d; offset:%d, value: %v\n", message.Partition, message.Offset, topup)

		account, err := m.repo.GetAccountByNorek(topup.Norek)
		if err != nil {
			return err
		}

		account.Saldo += topup.Gram

		trx := model.Transaction{
			Id:           shortid.MustGenerate(),
			Date:         time.Now().Unix(),
			Type:         "topup",
			RekeningId:   account.Id,
			Norek:        account.Norek,
			Gram:         topup.Gram,
			HargaTopup:   topup.HargaTopup,
			HargaBuyback: topup.HargaBuyback,
			Saldo:        account.Saldo,
		}

		logrus.Info(account)
		logrus.Info(topup)
		err = m.repo.InsertTransaksiAndUpdateBalance(trx, *account)
		if err != nil {
			return err
		}
	}

	return err
}
