package repository

import (
	"price-check-service/model"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type (
	PriceCheckRepo interface {
		GetLatestPrice() (*model.Price, error)
	}

	PriceCheckRepoImpl struct {
		db *sqlx.DB
	}
)

func NewPriceCheckRepo(db *sqlx.DB) PriceCheckRepo {
	return &PriceCheckRepoImpl{db: db}
}

const (
	selectLatestPriceQuery = `SELECT id, topup , buyback , "date" , admin_id FROM tbl_harga th ORDER BY "date" DESC LIMIT 1;`
)

func (r *PriceCheckRepoImpl) GetLatestPrice() (*model.Price, error) {
	var price model.Price
	err := r.db.Get(&price, selectLatestPriceQuery)
	if err != nil {
		logrus.Errorf("[price-check-service] Failed select data from database err: %s", err.Error())
		return nil, err
	}

	return &price, nil
}
