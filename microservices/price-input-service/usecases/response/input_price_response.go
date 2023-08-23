package response

type InputPriceResponse struct {
	Error   bool   `json:"error"`
	RefId   string `json:"reff_id,omitempty"`
	Message string `json:"message,omitempty"`
}
