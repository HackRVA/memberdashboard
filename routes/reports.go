package routes

import (
	"net/http"

	"github.com/HackRVA/memberserver/middleware/rbac"
)

type ReportsHTTPHandler interface {
	GetMemberCountsCharts(http.ResponseWriter, *http.Request)
	GetAccessStatsChart(http.ResponseWriter, *http.Request)
	GetMemberChurn(http.ResponseWriter, *http.Request)
}

func (r Router) setupReportsRoutes(reports ReportsHTTPHandler, accessControl rbac.AccessControl) {
	r.authedRouter.HandleFunc("/reports/membercounts", accessControl.Restrict(reports.GetMemberCountsCharts, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/reports/access", accessControl.Restrict(reports.GetAccessStatsChart, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/reports/churn", reports.GetMemberChurn)
}
