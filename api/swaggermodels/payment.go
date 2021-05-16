package api

import "memberserver/api/models"

// PaymentResponse response of payment chart information
// swagger:response getPaymentChartResponse
type getPaymentChartResponse struct {
	// in: body
	Body []models.PaymentChart
}

// swagger:parameters searchPaymentChartRequest
type searchPaymentChartRequest struct {
	// in:query
	Type string `json:"type"`
}
