package response

import "mutation-check-service/model"

type BuybackResponse struct {
	model.BaseResponse
	RefId        string `json:"reff_id,omitempty"`
	ResponseCode int    `json:"-"`
}
