package usecases

import (
	"encoding/json"
	"price-input-service/broker"
	"topup-service/usecases/response"

	"github.com/teris-io/shortid"
)

type (
	Topup struct {
		Gram         float32 `json:"gram"`
		HargaTopup   int64   `json:"harga_topup"`
		HargaBuyback int64   `json:"harga_buyback"`
		Norek        string  `json:"norek,omitempty"`
	}

	TopupUsecase interface {
		Topup(payload Topup) response.TopupResponse
	}

	TopupUsecaseImpl struct {
		messageBroker broker.MessageBroker
	}
)

func NewTopupUsecase(messageBroker broker.MessageBroker) TopupUsecase {
	return &TopupUsecaseImpl{
		messageBroker: messageBroker,
	}
}

func (p *TopupUsecaseImpl) Topup(payload Topup) response.TopupResponse {
	response := response.TopupResponse{}

	topic := "topup"
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
