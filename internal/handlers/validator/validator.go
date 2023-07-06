package validator

import (
	"errors"
	"github.com/fdjrn/dw-balance-history-service/internal/db/entity"
)

func ValidateRequest(payload interface{}) (interface{}, error) {
	var msg []string

	switch p := payload.(type) {
	case *entity.PaginatedRequest:
		if p.PartnerID == "" {
			msg = append(msg, "partnerId cannot be empty.")
		}

		if p.MerchantID == "" {
			msg = append(msg, "merchantId cannot be empty.")
		}

	default:

	}
	if len(msg) > 0 {
		return msg, errors.New("request validation status failed")
	}
	return msg, nil
}
