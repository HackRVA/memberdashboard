package api

import "memberserver/api/models"

// PaymentResponse response of payment chart information
// swagger:response getPaymentChartResponse
type getPaymentChartResponse struct {
	// in: body
	Body []models.PaymentChart
}
