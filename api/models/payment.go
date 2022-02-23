package models

import (
	"time"

	"github.com/Rhymond/go-money"
)

// PaymentProvider enum
type PaymentProvider int

const (
	//QuickBooks ... payment provider
	QuickBooks PaymentProvider = iota
	//Paypal ... payment provider
	Paypal
)

// ChartOptions -- config option for the chart
type ChartOptions struct {
	Title     string  `json:"title"`
	CurveType string  `json:"curveType"`
	PieHole   float64 `json:"pieHole"`
	Legend    string  `json:"legend"`
}

// ChartRow -- chart row info
type ChartRow struct {
	MemberLevel, MemberCount interface{}
}

// ChartCol -- chart col info
type ChartCol struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

// ReportChart - deliver information to the frontend so that
//   we can display a monthly payment chart
type ReportChart struct {
	Options ChartOptions    `json:"options"`
	Type    string          `json:"type"`
	Rows    [][]interface{} `json:"rows"`
	Cols    []ChartCol      `json:"cols"`
}

// Payment represents a payment made
// this will be pulled down from the providers
type Payment struct {
	ID string
	// Date is when the payment was made
	Date     time.Time
	Amount   money.Money
	Provider PaymentProvider
	MemberID string
	Email    string
	Name     string
}

// PastDueAccount represents accounts that do not have a recent payment recorded
type PastDueAccount struct {
	MemberId             string
	Name                 string
	Email                string
	LastPaymentDate      time.Time
	DaysSinceLastPayment int
}
