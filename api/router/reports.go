package router

import (
	"memberserver/api/rbac"
	"net/http"
)

type ReportsHTTPHandler interface {
	GetMemberCounts(http.ResponseWriter, *http.Request)
}

func (r Router) setupReportsRoutes(reports ReportsHTTPHandler, accessControl rbac.AccessControl) {
	// swagger:route GET /api/member/stats stats searchPaymentChartRequest
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
	//
	//     Responses:
	//       200: getPaymentChartResponse
	r.authedRouter.HandleFunc("/reports/membercounts", accessControl.Restrict(reports.GetMemberCounts, []rbac.UserRole{rbac.Admin}))
}
