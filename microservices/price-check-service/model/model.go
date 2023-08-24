package model

type (
	Price struct {
		Id      string `db:"id"`
		Topup   int64  `db:"topup"`
		Buyback int64  `db:"buyback"`
		Date    int64  `db:"date"`
		AdminId string `db:"admin_id"`
	}
)
