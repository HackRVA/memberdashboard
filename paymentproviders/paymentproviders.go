package paymentproviders

import "time"

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
	ID         string
	Date       time.Time
	Amount     float32
	Provider   PaymentProvider
	CustomerID uint8
	UpdatedAt  time.Time
}

// GetPayments reach out the payment providers and download
// payments
func GetPayments() {
	getPaypalPayments()
	getQBPayments()
}

func getPaypalPayments() {

}

func getQBPayments() {
}

func hashPayment() string {
	var hash string

	return hash
}
