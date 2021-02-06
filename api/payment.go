package api

import (
	"encoding/json"
	"memberserver/database"
	"net/http"
)

type chartOptions struct {
	Title     string  `json:"title"`
	CurveType string  `json:"curveType"`
	PieHole   float64 `json:"pieHole"`
	Legend    string  `json:"legend"`
}

type chartRow struct {
	MemberLevel, MemberCount interface{}
}

type chartCol struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

// PaymentChart - deliver information to the frontend so that
//   we can display a monthly payment chart
type PaymentChart struct {
	Options chartOptions    `json:"options"`
	Type    string          `json:"type"`
	Rows    [][]interface{} `json:"rows"`
	Cols    []chartCol      `json:"cols"`
}

// PaymentResponse response of payment chart information
// swagger:response getPaymentChartResponse
type getPaymentChartResponse struct {
	PaymentCharts []PaymentChart `json:"paymentCharts"`
}

// countMemberLevels take in a list of payments and return
//   formatted data to be used in payment charts
func countMemberLevels(payments []int64) map[database.MemberLevel]uint8 {
	counts := make(map[database.MemberLevel]uint8)

	// set counts to 0
	//  kind of cheating here by setting them directly
	counts[database.Inactive] = 0
	counts[database.Classic] = 0
	counts[database.Standard] = 0
	counts[database.Premium] = 0

	for _, p := range payments {
		if v, found := database.MemberLevelFromAmount[p]; found {
			counts[v]++
		}
	}

	return counts
}

func makeMemberCountTrendChart(payments map[string][]int64) PaymentChart {
	var pc PaymentChart
	pc.Options.Title = "Membership Trends"
	pc.Type = "line"
	pc.Options.CurveType = "function"
	pc.Options.Legend = "bottom"
	pc.Cols = []chartCol{{Label: "Month", Type: "string"}, {Label: "MemberCount", Type: "number"}}

	for k, p := range payments {
		var row []interface{}
		row = append(row, k)
		row = append(row, len(p))
		pc.Rows = append(pc.Rows, row)
	}
	return pc
}

func (a API) getPaymentChart(w http.ResponseWriter, req *http.Request) {
	paymentList, err := a.db.GetPayments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paymentMapByDate := make(map[string][]int64)

	// get the rows together
	for _, p := range paymentList {
		paymentMapByDate[p.Date.Format("Jan-06")] = append(paymentMapByDate[p.Date.Format("Jan-06")], int64(p.Amount.AsMajorUnits()))
	}

	// now the rows are together, but they are in the form of a map
	// let's massage it out to match our chart contract

	var paymentCharts []PaymentChart

	paymentCharts = append(paymentCharts, makeMemberCountTrendChart(paymentMapByDate))

	for k, paymentsByMonth := range paymentMapByDate {
		var pc PaymentChart
		pc.Options.Title = k + " - Membership Distribution"
		pc.Options.PieHole = 0.4
		pc.Type = "pie"

		pc.Cols = []chartCol{{Label: "Month", Type: "string"}, {Label: "MemberLevelCount", Type: "number"}}
		levels := countMemberLevels(paymentsByMonth)

		for level, count := range levels {
			var row []interface{}
			row = append(row, (string)(database.MemberLevelToStr[level]))
			row = append(row, int(count))
			pc.Rows = append(pc.Rows, row)
		}

		paymentCharts = append(paymentCharts, pc)
	}

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(paymentCharts)
	w.Write(j)
}
