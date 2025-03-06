package swagger

import "github.com/HackRVA/memberserver/models"

// PaymentResponse response of payment chart information
// swagger:response getPaymentChartResponse
type getPaymentChartResponse struct {
	// in: body
	Body []models.ReportChart
}

// swagger:parameters membershipChartRequest
type membershipChartRequest struct {
	// in:query
	Type string `json:"type"`
}

// PaymentResponse response of payment chart information
// swagger:response getAccessStatsChartResponse
type getAccessStatsChartResponse struct {
	// in: body
	Body []models.AccessStats
}

// PaymentResponse response of payment chart information
// swagger:response getMemberChurnResponse
type getMemberChurnResponse struct {
	// in: body
	Body []models.AccessStats
}

// swagger:parameters accessChartRequest
type accessChartRequest struct {
	// in:query
	ResourceName string `json:"resourceName"`
}
