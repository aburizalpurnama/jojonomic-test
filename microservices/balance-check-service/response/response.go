package response

type BalanceCheckResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    *Data  `json:"data,omitempty"`
}

type Data struct {
	Norek string  `json:"norek"`
	Saldo float32 `json:"saldo"`
}
