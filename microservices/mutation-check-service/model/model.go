package model

type (
	MutationCheckRequest struct {
		Norek     string `json:"norek"`
		StartDate int64  `json:"start_date"`
		EndDate   int64  `json:"end_date"`
	}

	MutationCheckResponse struct {
		BaseResponse
		Data []Transcation `json:"data,omitempty"`
	}

	Transcation struct {
		Date         int64   `json:"date"`
		Type         string  `json:"type"`
		Gram         float32 `json:"gram"`
		HargaTopup   int64   `json:"harga_topup"`
		HargaBuyback int64   `json:"harga_buyback"`
		Saldo        float32 `json:"saldo"`
	}

	BaseResponse struct {
		Error   bool   `json:"error"`
		Message string `json:"message,omitempty"`
	}
)
