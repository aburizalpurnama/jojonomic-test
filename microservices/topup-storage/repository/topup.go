package repository

import (
	"database/sql"
	"topup-storage/model"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type (
	TopupRepo interface {
		InsertTransaksiAndUpdateBalance(transaction model.Transaction, account model.Account) error
		UpdateBalance(data model.Account) error
		GetAccountByNorek(norek string) (*model.Account, error)
	}

	TopupRepoImpl struct {
		db *sqlx.DB
		tx *sql.Tx
	}
)

func NewTopupRepo(db *sqlx.DB) TopupRepo {
	return &TopupRepoImpl{db: db}
}

const (
	inputTransactionQuery = `
	INSERT INTO tbl_transaksi (id, "date" , "type", rekening_id, norek, gram, harga_topup, harga_buyback, saldo)
	VALUES (:id, :date , :type, :rekening_id, :norek, :gram, :harga_topup, :harga_buyback, :saldo);`
	updateBalanceAccountQuery = `UPDATE tbl_rekening SET saldo = $1 WHERE id = $2;`
	getAccountQuery           = `SELECT id , norek , saldo ,updated_at FROM tbl_rekening WHERE norek = $1;`
)

func (r *TopupRepoImpl) InsertTransaksiAndUpdateBalance(transaction model.Transaction, account model.Account) error {
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(inputTransactionQuery, transaction)
	if err != nil {
		tx.Rollback()
		logrus.Errorf("[topup-storage] Failed insert transaction err: %s", err.Error())
		return err
	}

	_, err = tx.Exec(updateBalanceAccountQuery, account.Saldo, account.Id)
	if err != nil {
		tx.Rollback()
		logrus.Errorf("[topup-storage] Failed update balance err: %s", err.Error())
		return err
	}

	err = tx.Commit()

	return err
}

func (r *TopupRepoImpl) UpdateBalance(data model.Account) error {
	_, err := r.db.Exec(updateBalanceAccountQuery, data.Saldo, data.Id)
	if err != nil {
		logrus.Errorf("[topup-storage] Failed update balance err: %s", err.Error())
	}

	return err
}

func (r *TopupRepoImpl) GetAccountByNorek(norek string) (*model.Account, error) {
	var account model.Account
	err := r.db.Get(&account, getAccountQuery, norek)
	if err != nil {
		logrus.Errorf("[topup-storage] Failed get account by norek err: %s", err.Error())
		return nil, err
	}

	return &account, nil
}
