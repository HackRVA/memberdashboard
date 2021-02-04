package paymentproviders

import (
	"time"

	"github.com/Rhymond/go-money"
	log "github.com/sirupsen/logrus"
)

// PaymentProvider enum
type PaymentProvider int

const (
	//QuickBooks ... payment provider
	QuickBooks PaymentProvider = iota
	//Paypal ... payment provider
	Paypal
)

type paypalAuth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// Payment represents a payment made
// this will be pulled down from the providers
type Payment struct {
	ID string
	// Date is when the payment was made
	Date       time.Time
	Amount     money.Money
	Provider   PaymentProvider
	CustomerID uint8
	Email      string
}

// GetPayments reach out the payment providers and download
// payments
func GetPayments() {
	startDate := time.Now().AddDate(0, -1, 0).Format(time.RFC3339) // subtract one month
	endDate := time.Now().Format(time.RFC3339)

	payments, err := getPaypalPayments(startDate, endDate)
	if err != nil {
		log.Errorf("error getting payments: %s\n", err.Error())
	}

	for _, p := range payments {
		log.Printf("%s, %s, %s", p.Email, p.Amount.Display(), p.Date.String())
	}

	// TODO: add payments to DB
}

func hashPayment() string {
	var hash string

	return hash
}
