package usecases

import (
	"price-check-service/repository"
	"price-check-service/usecases/response"
)

type (
	PriceCheckUsecase interface {
		CheckPrice() (*response.CheckPriceResponse, error)
	}

	PriceCheckUsecaseImpl struct {
		repo repository.PriceCheckRepo
	}
)

func NewPriceCheckUsecase(repo repository.PriceCheckRepo) PriceCheckUsecase {
	return &PriceCheckUsecaseImpl{
		repo: repo,
	}
}

func (p *PriceCheckUsecaseImpl) CheckPrice() (*response.CheckPriceResponse, error) {
	var resp response.CheckPriceResponse
	latestPrice, err := p.repo.GetLatestPrice()
	if err != nil {
		resp.Error = true
		resp.Message = err.Error()
		return &resp, err
	}

	resp.Data = response.CheckPriceData{
		Topup:   latestPrice.Topup,
		Buyback: latestPrice.Buyback,
	}
	resp.Error = false

	return &resp, nil
}
