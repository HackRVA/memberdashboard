package paymentproviders

import (
	"time"
)

// PaymentProvider enum
type PaymentProvider int

const (
	//QuickBooks ... payment provider
	QuickBooks PaymentProvider = iota
	//Paypal ... payment provider
	Paypal
)

// Payment represents a payment made
// this will be pulled down from the providers
type Payment struct {
	ID string
	// Date is when the payment was made
	Date       time.Time
	Amount     float32
	Provider   PaymentProvider
	CustomerID uint8
	// UpdatedAt is when the record was updated in our DB
	UpdatedAt time.Time
}

// GetPayments reach out the payment providers and download
// payments
func GetPayments() {
	startDate, _ := time.Parse(time.RFC3339, "2014-07-01T00:00:00-0700")
	endDate, _ := time.Parse(time.RFC3339, "2014-07-30T23:59:59-0700")
	getPaypalPayments(startDate, endDate)
	getQBPayments()
}

func getQBPayments() {
}

func hashPayment() string {
	var hash string

	return hash
}
