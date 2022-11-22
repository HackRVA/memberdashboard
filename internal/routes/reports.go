package routes

import (
	"net/http"

	"github.com/HackRVA/memberserver/internal/middleware/rbac"
)

type ReportsHTTPHandler interface {
	GetMemberCountsCharts(http.ResponseWriter, *http.Request)
	GetAccessStatsChart(http.ResponseWriter, *http.Request)
	GetMemberChurn(http.ResponseWriter, *http.Request)
}

func (r Router) setupReportsRoutes(reports ReportsHTTPHandler, accessControl rbac.AccessControl) {
	// swagger:route GET /api/reports/membercounts stats membershipChartRequest
	//
	// Get Chart information of monthly member counts
	//
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//     - basicAuth:
	//
	//     Responses:
	//       200: getPaymentChartResponse
	r.authedRouter.HandleFunc("/reports/membercounts", accessControl.Restrict(reports.GetMemberCountsCharts, []rbac.UserRole{rbac.Admin}))
	// swagger:route GET /api/reports/access stats accessChartRequest
	//
	// Get Chart information of daily access event counts
	//
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//     - basicAuth:
	//
	//     Responses:
	//       200: getAccessStatsChartResponse
	r.authedRouter.HandleFunc("/reports/access", accessControl.Restrict(reports.GetAccessStatsChart, []rbac.UserRole{rbac.Admin}))
	// swagger:route GET /api/reports/churn stats
	//
	// Get Chart information of daily access event counts
	//
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//     - basicAuth:
	//
	//     Responses:
	//       200: getMemberChurnResponse
	r.authedRouter.HandleFunc("/reports/churn", reports.GetMemberChurn)
}
