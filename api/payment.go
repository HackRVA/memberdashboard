package api

import (
	"encoding/json"
	"memberserver/api/models"
	"memberserver/database"
	"net/http"

	"github.com/emirpasic/gods/maps/linkedhashmap"
)

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

func makeMemberCountTrendChart(payments linkedhashmap.Map) models.PaymentChart {
	var pc models.PaymentChart
	pc.Options.Title = "Membership Trends"
	pc.Type = "line"
	pc.Options.CurveType = "function"
	pc.Options.Legend = "bottom"
	pc.Cols = []models.ChartCol{{Label: "Month", Type: "string"}, {Label: "Member Count", Type: "number"}}

	it := payments.Iterator()

	for it.Next() {
		var row []interface{}
		row = append(row, it.Key())
		var paymentAmounts []int64 = it.Value().([]int64)
		row = append(row, len(paymentAmounts))
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

	paymentMapByDate := linkedhashmap.New()

	// get the rows together
	for _, p := range paymentList {
		_, found := paymentMapByDate.Get(p.Date.Format("Jan-06"))
		var paymentAmounts []int64

		if found {
			monthPaymentAmounts, _ := paymentMapByDate.Get(p.Date.Format("Jan-06"))
			paymentAmounts = monthPaymentAmounts.([]int64)
		}

		paymentAmounts = append(paymentAmounts, int64(p.Amount.AsMajorUnits()))
		paymentMapByDate.Put(p.Date.Format(("Jan-06")), paymentAmounts)
	}

	// now the rows are together, but they are in the form of a map
	// let's massage it out to match our chart contract

	var paymentCharts []models.PaymentChart

	paymentCharts = append(paymentCharts, makeMemberCountTrendChart(*paymentMapByDate))

	it := paymentMapByDate.Iterator()

	for it.Next() {
		var pc models.PaymentChart
		pc.Options.Title = it.Key().(string) + " - Membership Distribution"
		pc.Options.PieHole = 0.4
		pc.Type = "pie"

		pc.Cols = []models.ChartCol{{Label: "Month", Type: "string"}, {Label: "MemberLevelCount", Type: "number"}}
		levels := countMemberLevels(it.Value().([]int64))

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
