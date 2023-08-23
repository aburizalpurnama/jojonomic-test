package request

type UpdatePriceRequest struct {
	AdminId string `json:"admin_id"`
	Topup   int64  `json:"harga_topup"`
	Buyback int64  `json:"harga_buyback"`
}
