package model

type (
	Price struct {
		Id      string `db:"id" json:"-"`
		Topup   int64  `db:"topup" json:"harga_topup"`
		Buyback int64  `db:"buyback" json:"harga_buyback"`
		Date    int64  `db:"date" json:"-"`
		AdminId string `db:"admin_id" json:"admin_id"`
	}
)
