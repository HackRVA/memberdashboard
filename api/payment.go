package api

import (
	"encoding/json"
	"memberserver/api/models"
	"memberserver/payments"

	"memberserver/payments/listener"
	"net/http"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	log "github.com/sirupsen/logrus"
)

// countMemberLevels take in a list of payments and return
//   formatted data to be used in payment charts
func countMemberLevels(payments []int64) map[models.MemberLevel]uint8 {
	counts := make(map[models.MemberLevel]uint8)

	// set counts to 0
	//  kind of cheating here by setting them directly
	counts[models.Inactive] = 0
	counts[models.Classic] = 0
	counts[models.Standard] = 0
	counts[models.Premium] = 0

	for _, p := range payments {
		if v, found := models.MemberLevelFromAmount[p]; found {
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
	chartType := req.URL.Query().Get("type")
	var paymentCharts []models.PaymentChart

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

	if len(chartType) > 0 {
		if chartType == "line" {
			paymentCharts = append(paymentCharts, makeMemberCountTrendChart(*paymentMapByDate))
		}

		if chartType == "pie" {
			paymentCharts = makeMemberDistributionChart(*paymentMapByDate, paymentCharts)
		}
	}

	if len(chartType) == 0 {
		paymentCharts = append(paymentCharts, makeMemberCountTrendChart(*paymentMapByDate))
		paymentCharts = makeMemberDistributionChart(*paymentMapByDate, paymentCharts)
	}

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(paymentCharts)
	w.Write(j)
}

func makeMemberDistributionChart(payments linkedhashmap.Map, paymentCharts []models.PaymentChart) []models.PaymentChart {
	// now the rows are together, but they are in the form of a map
	// let's massage it out to match our chart contract
	it := payments.Iterator()

	for it.Next() {
		var pc models.PaymentChart
		pc.Options.Title = it.Key().(string)
		pc.Options.PieHole = 0.4
		pc.Type = "pie"

		pc.Cols = []models.ChartCol{{Label: "Month", Type: "string"}, {Label: "MemberLevelCount", Type: "number"}}
		levels := countMemberLevels(it.Value().([]int64))

		for level, count := range levels {
			var row []interface{}
			row = append(row, (string)(models.MemberLevelToStr[level]))
			row = append(row, int(count))
			pc.Rows = append(pc.Rows, row)
		}

		paymentCharts = append(paymentCharts, pc)
	}

	return paymentCharts
}

// PaypalSubscriptionWebHookHandler paypal will tell us when a new subscription is created.
//   We can use this to add a member to our database.  We don't have to give them
//   access to anything at this time, but it will make it easier to assign them an RFID fob
func (api *API) PaypalSubscriptionWebHookHandler(err error, n *listener.PaypalNotification) {
	if err != nil {
		log.Printf("IPN error: %v", err)
		return
	}

	log.Printf("event type: %s", n.EventType)
	log.Printf("event resource type: %s", n.ResourceType)
	log.Printf("summary: %s", n.Summary)
	log.Printf("name: %s", n.Resource.Subscriber.Name.GivenName+" "+n.Resource.Subscriber.Name.SurName)

	if n.EventType != "BILLING.SUBSCRIPTION.CREATED" {
		return
	}

	paymentProvider := payments.Setup(api.db)

	newMember, err := paymentProvider.GetSubscription(n.ID)
	if err != nil {
		log.Errorf("error parsing member from webhook: %v", err)
	}

	api.db.ProcessMember(newMember)
}
