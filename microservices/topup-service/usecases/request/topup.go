package request

type TopupRequest struct {
	Gram  string `json:"gram"`
	Harga string `json:"harga"`
	Norek string `json:"norek,omitempty"`
}
