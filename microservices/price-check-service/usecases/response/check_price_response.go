package response

type (
	CheckPriceResponse struct {
		Error   bool           `json:"error"`
		Data    CheckPriceData `json:"data,omitempty"`
		Message string         `json:"message,omitempty"`
	}

	CheckPriceData struct {
		Topup   int64 `json:"harga_topup"`
		Buyback int64 `json:"harga_buyback"`
	}
)
