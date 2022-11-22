package controllers

import (
	"net/http"
	"time"

	"github.com/HackRVA/memberserver/internal/models"
	"github.com/HackRVA/memberserver/internal/services/report"
)

type ReportsServer struct {
	service report.ReportService
	Logger  Logger
}

func (r *ReportsServer) GetAccessStatsChart(w http.ResponseWriter, req *http.Request) {
	resourceName := req.URL.Query().Get("resourceName")
	day := req.URL.Query().Get("day")

	var d time.Time
	var err error

	if len(resourceName) == 0 {
		preconditionFailed(w, "invalid resourceName")
		return
	}

	if len(day) > 0 {
		d, err = time.Parse("", day)
		if err != nil {
			r.Logger.Errorf("error parsing time")
		}
	}

	charts, err := r.service.GetAccessStatsChart(d, resourceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok(w, charts)
}

func (r *ReportsServer) GetMemberCountsCharts(w http.ResponseWriter, req *http.Request) {
	chartType := req.URL.Query().Get("type")
	month := req.URL.Query().Get("month")

	if len(month) > 0 {
		date, err := time.Parse("", month)
		if err != nil {
			preconditionFailed(w, "invalid month")
			return
		}

		ok(w, r.service.GetMemberCountsChartByMonth(date))
		return
	}

	charts, err := r.service.GetMemberCountsCharts(chartType)
	if err != nil {
		return
	}
	ok(w, charts)
}

func (r *ReportsServer) GetMemberChurn(w http.ResponseWriter, req *http.Request) {
	churn, err := r.service.GetMemberChurn()
	if err != nil {
		internalServerError(w, "error getting member churn")
		return
	}

	ok(w, models.MemberChurn{
		Churn: churn,
	})
}
