package repository

import (
	"topup-storage/model"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type (
	MutationCheckRepo interface {
		SelectTransactionsByNorek(norek string, start int64, end int64) ([]model.Transaction, error)
	}

	MutationCheckRepoImpl struct {
		db *sqlx.DB
	}
)

func NewMutationCheckRepo(db *sqlx.DB) MutationCheckRepo {
	return &MutationCheckRepoImpl{db: db}
}

const (
	selectTransactionsByNorekQuery = `SELECT id , "date" , "type" , rekening_id , norek , gram , harga_topup , harga_buyback , saldo FROM tbl_transaksi WHERE norek = $1 AND "date" BETWEEN $2 AND $3 ORDER BY "date" DESC;`
)

func (r *MutationCheckRepoImpl) SelectTransactionsByNorek(norek string, start int64, end int64) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Select(&transactions, selectTransactionsByNorekQuery, norek, start, end)
	if err != nil {
		logrus.Errorf("[topup-storage] Failed get account by norek err: %s", err.Error())
		return nil, err
	}

	return transactions, nil
}
