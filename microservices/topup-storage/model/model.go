package model

type (
	Transaction struct {
		Id           string  `db:"id"`
		Date         int64   `db:"date"`
		Type         string  `db:"type"`
		RekeningId   string  `db:"rekening_id"`
		Norek        string  `db:"norek"`
		Gram         float32 `db:"gram"`
		HargaTopup   int64   `db:"harga_topup"`
		HargaBuyback int64   `db:"harga_buyback"`
		Saldo        float32 `db:"saldo"`
	}

	Account struct {
		Id        string  `db:"id"`
		Norek     string  `db:"norek"`
		Saldo     float32 `db:"saldo"`
		UpdatedAt int64   `db:"updated_at"`
	}
)
