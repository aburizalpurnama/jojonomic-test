package request

type BuybackRequest struct {
	Gram  float32 `json:"gram"`
	Harga int64   `json:"harga"`
	Norek string  `json:"norek"`
}

type BuybackPayload struct {
	Transaction Transaction `json:"transaction"`
	Account     Account     `json:"account"`
}

type Transaction struct {
	Id           string  `json:"id"`
	Date         int64   `json:"date"`
	Type         string  `json:"type"`
	RekeningId   string  `json:"rekening_id"`
	Norek        string  `json:"norek"`
	Gram         float32 `json:"gram"`
	HargaTopup   int64   `json:"harga_topup"`
	HargaBuyback int64   `json:"harga_buyback"`
	Saldo        float32 `json:"saldo"`
}

type Account struct {
	Id        string  `json:"id"`
	Norek     string  `json:"norek"`
	Saldo     float32 `json:"saldo"`
	UpdatedAt int64   `json:"updated_at"`
}
