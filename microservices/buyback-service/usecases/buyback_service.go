package usecases

import (
	"buyback-service/usecases/request"
	"buyback-service/usecases/response"
	"database/sql"
	"encoding/json"
	"net/http"
	"price-check-service/repository"
	"price-input-service/broker"
	"time"
	"topup-service/usecases"
	topupRepo "topup-storage/repository"

	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

type (
	BuybackUsecase interface {
		Buyback(payload request.BuybackRequest) (response.BuybackResponse, error)
	}

	BuybackUsecaseImpl struct {
		messageBroker  broker.MessageBroker
		topupRepo      topupRepo.TopupRepo
		priceCheckRepo repository.PriceCheckRepo
	}
)

func NewBuybackUsecase(messageBroker broker.MessageBroker, topupRepo topupRepo.TopupRepo, priceCheckRepo repository.PriceCheckRepo) BuybackUsecase {
	return &BuybackUsecaseImpl{
		messageBroker:  messageBroker,
		topupRepo:      topupRepo,
		priceCheckRepo: priceCheckRepo,
	}
}

func (p *BuybackUsecaseImpl) Buyback(payload request.BuybackRequest) (response.BuybackResponse, error) {
	response := response.BuybackResponse{}

	account, err := p.topupRepo.GetAccountByNorek(payload.Norek)
	switch err {
	default:
		response.Error = true
		response.Message = err.Error()
		response.ResponseCode = http.StatusInternalServerError
		return response, err
	case sql.ErrNoRows:
		response.Error = true
		response.Message = "account not found"
		response.ResponseCode = http.StatusNotFound
		return response, err
	case nil:
	}

	price, err := p.priceCheckRepo.GetLatestPrice()
	if err != nil {
		response.Error = true
		response.Message = err.Error()
		response.ResponseCode = http.StatusInternalServerError
		return response, err
	}

	if price.Buyback != payload.Harga {
		response.Error = true
		response.Message = "invalid price"
		response.ResponseCode = http.StatusBadRequest
		return response, err
	}

	count := usecases.CountDecimalPlaces(payload.Gram)
	if count > 3 {
		logrus.Info("invalid gram, max 3 decimal places")
		response.Error = true
		response.Message = "invalid gram, max 3 decimal places"
		response.ResponseCode = http.StatusBadRequest
		return response, err
	}

	if account.Saldo < payload.Gram {
		response.Error = true
		response.Message = "not enough balance"
		response.ResponseCode = http.StatusBadRequest
		return response, err
	}

	topic := "buyback"
	key := shortid.MustGenerate()

	buybackPayload := request.BuybackPayload{
		Transaction: request.Transaction{
			Id:           key,
			Date:         time.Now().Unix(),
			Type:         "buyback",
			RekeningId:   account.Id,
			Norek:        account.Norek,
			Gram:         payload.Gram,
			HargaTopup:   price.Topup,
			HargaBuyback: price.Buyback,
			Saldo:        account.Saldo - payload.Gram,
		},
		Account: request.Account{
			Id:    account.Id,
			Saldo: account.Saldo - payload.Gram,
		},
	}

	value, _ := json.Marshal(buybackPayload)
	err = p.messageBroker.ProcudeData(topic, key, value)
	if err != nil {
		response.Error = true
		response.Message = err.Error()
		response.RefId = key
		response.ResponseCode = http.StatusInternalServerError

		return response, err
	}

	response.Error = false
	response.RefId = key
	response.ResponseCode = http.StatusOK

	return response, err
}
