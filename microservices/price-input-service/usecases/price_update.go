package usecases

import (
	"encoding/json"
	"price-input-service/broker"
	"price-input-service/usecases/request"
	"price-input-service/usecases/response"

	"github.com/teris-io/shortid"
)

type (
	PriceInputUsecase interface {
		InputPrice(payload request.UpdatePriceRequest) response.InputPriceResponse
	}

	priceInputUsecaseImpl struct {
		messageBroker broker.MessageBroker
	}
)

func NewPriceInputUsecase(messageBroker broker.MessageBroker) PriceInputUsecase {
	return &priceInputUsecaseImpl{
		messageBroker: messageBroker,
	}
}

func (p *priceInputUsecaseImpl) InputPrice(payload request.UpdatePriceRequest) response.InputPriceResponse {
	response := response.InputPriceResponse{}

	topic := "input-harga"
	key := shortid.MustGenerate()
	value, _ := json.Marshal(payload)

	err := p.messageBroker.ProcudeData(topic, key, value)
	if err != nil {
		response.Error = true
		response.Message = err.Error()
		response.RefId = key

		return response
	}

	response.Error = false
	response.RefId = key

	return response
}
