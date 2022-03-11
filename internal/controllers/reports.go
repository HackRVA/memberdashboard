package controllers

import (
	"memberserver/internal/datastore"
	"memberserver/internal/models"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type ReportsServer struct {
	store datastore.ReportStore
}

func (r *ReportsServer) GetMemberCounts(w http.ResponseWriter, req *http.Request) {
	chartType := req.URL.Query().Get("type")
	month := req.URL.Query().Get("month")

	if len(month) > 0 {
		date, err := time.Parse("", month)
		if err != nil {
			http.Error(w, "error looking up counts by month - use a valid date", http.StatusNotFound)
			return
		}
		ok(w, makeDistritutionChartByMonth(date, r.store))
		return
	}

	var charts []models.ReportChart
	memberCounts, err := r.store.GetMemberCounts()
	if err != nil {
		http.Error(w, "error getting member counts", http.StatusNotFound)
		return
	}

	if len(chartType) > 0 {
		if chartType == "line" {
			charts = append(charts, makeMemberTrendChart(memberCounts))
		}

		if chartType == "pie" {
			charts = makeMemberDistributionChart(memberCounts)
		}
	}

	if len(chartType) == 0 {
		charts = append(charts, makeMemberTrendChart(memberCounts))
		charts = append(charts, makeMemberDistributionChart(memberCounts)...)
	}

	ok(w, charts)
}

func makeMemberTrendChart(counts []models.MemberCount) models.ReportChart {
	var chart models.ReportChart
	chart.Options.Title = "Membership Trends"
	chart.Type = "line"
	chart.Options.CurveType = "function"
	chart.Options.Legend = "bottom"
	chart.Cols = []models.ChartCol{{Label: "Month", Type: "string"}, {Label: "Member Count", Type: "number"}}

	for _, monthCount := range counts {
		var row []interface{}
		row = append(row, monthCount.Month.Format("Jan-06"))
		// explicitly exclude credited
		row = append(row, monthCount.Classic+monthCount.Standard+monthCount.Premium)
		chart.Rows = append(chart.Rows, row)
	}
	return chart
}

func makeDistritutionChartByMonth(month time.Time, store datastore.ReportStore) models.ReportChart {
	var distributionChart models.ReportChart
	memberCount, err := store.GetMemberCountByMonth(month)
	if err != nil {
		log.Errorf("error getting member counts")
		return distributionChart
	}

	var chart models.ReportChart
	chart.Options.Title = month.Format("Jan-06")
	chart.Options.PieHole = 0.4
	chart.Type = "pie"

	chart.Cols = []models.ChartCol{{Label: "Month", Type: "string"}, {Label: "MemberLevelCount", Type: "number"}}

	levels := make(map[models.MemberLevel]uint8)
	levels[models.Credited] = uint8(memberCount.Credited)
	levels[models.Classic] = uint8(memberCount.Classic)
	levels[models.Standard] = uint8(memberCount.Standard)
	levels[models.Premium] = uint8(memberCount.Premium)

	for level, count := range levels {
		var row []interface{}
		row = append(row, (string)(models.MemberLevelToStr[level]))
		row = append(row, int(count))
		chart.Rows = append(chart.Rows, row)
	}

	return distributionChart
}

func makeMemberDistributionChart(counts []models.MemberCount) []models.ReportChart {
	var distributionCharts []models.ReportChart
	for _, monthCount := range counts {
		var chart models.ReportChart
		chart.Options.Title = monthCount.Month.Format("Jan-06")
		chart.Options.PieHole = 0.4
		chart.Type = "pie"

		chart.Cols = []models.ChartCol{{Label: "Month", Type: "string"}, {Label: "MemberLevelCount", Type: "number"}}
		levels := make(map[models.MemberLevel]uint8)
		levels[models.Credited] = uint8(monthCount.Credited)
		levels[models.Classic] = uint8(monthCount.Classic)
		levels[models.Standard] = uint8(monthCount.Standard)
		levels[models.Premium] = uint8(monthCount.Premium)

		for level, count := range levels {
			var row []interface{}
			row = append(row, (string)(models.MemberLevelToStr[level]))
			row = append(row, int(count))
			chart.Rows = append(chart.Rows, row)
		}

		distributionCharts = append(distributionCharts, chart)
	}
	return distributionCharts
}
