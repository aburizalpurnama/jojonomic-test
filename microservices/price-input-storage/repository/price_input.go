package repository

import (
	"price-input-storage/model"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type (
	PriceInputRepo interface {
		Insert(data model.Price) error
	}

	PriceInputRepoImpl struct {
		db *sqlx.DB
	}
)

func NewPriceInputRepo(db *sqlx.DB) PriceInputRepo {
	return &PriceInputRepoImpl{db: db}
}

const (
	inputPriceQuery = `
	INSERT INTO tbl_harga (id, topup, buyback, date, admin_id) 
	VALUES (:id, :topup, :buyback, :date, :admin_id)`
)

func (r *PriceInputRepoImpl) Insert(data model.Price) error {
	_, err := r.db.NamedExec(inputPriceQuery, data)
	if err != nil {
		logrus.Errorf("[price-input-storage] Failed insert to database err: %s", err.Error())
	}

	return err
}
