package report

import (
	"errors"
	"fmt"
	"time"

	"github.com/HackRVA/memberserver/internal/datastore"
	"github.com/HackRVA/memberserver/internal/models"

	"github.com/sirupsen/logrus"
)

type ReportService interface {
	GetAccessStatsChart(date time.Time, resourceName string) (models.ReportChart, error)
	GetMemberChurn() (int, error)
	GetMemberCountsCharts(chartType string) ([]models.ReportChart, error)
	GetMemberCountsChartByMonth(date time.Time) models.ReportChart
}

type Report struct {
	Store datastore.ReportStore
}

var (
	errNotFound = errors.New("not found")
)

func (r Report) GetAccessStatsChart(date time.Time, resourceName string) (models.ReportChart, error) {
	accessStats, err := r.Store.GetAccessStats(date, resourceName)
	if err != nil {
		return models.ReportChart{}, err
	}

	return makeAccessTrendChart(accessStats, resourceName), nil
}

func (r Report) GetMemberChurn() (int, error) {
	return r.Store.GetMemberChurn()
}

func (r Report) GetMemberCountsChartByMonth(date time.Time) models.ReportChart {
	return makeDistritutionChartByMonth(date, r.Store)
}

func (r Report) GetMemberCountsCharts(chartType string) ([]models.ReportChart, error) {
	var charts []models.ReportChart
	memberCounts, err := r.Store.GetMemberCounts()
	if err != nil {
		return []models.ReportChart{}, errNotFound
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

	return charts, nil
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
		logrus.Errorf("error getting member counts")
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

func makeAccessTrendChart(stats []models.AccessStats, resourceName string) models.ReportChart {
	var chart models.ReportChart
	chart.Options.Title = fmt.Sprintf("%s Access Trends", resourceName)
	chart.Type = "line"
	chart.Options.CurveType = "function"
	chart.Options.Legend = "bottom"
	chart.Cols = []models.ChartCol{{Label: "Day", Type: "string"}, {Label: "Access Count", Type: "number"}}

	for _, s := range stats {
		var row []interface{}
		row = append(row, s.Date.Format("Jan-06-19"))
		// explicitly exclude credited
		row = append(row, s.AccessCount)
		chart.Rows = append(chart.Rows, row)
	}
	return chart
}
