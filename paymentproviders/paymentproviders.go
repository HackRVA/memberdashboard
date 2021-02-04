package paymentproviders

import (
	"memberserver/database"
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
	Date     time.Time
	Amount   money.Money
	Provider PaymentProvider
	MemberID string
	Email    string
	Name     string
}

// GetPayments reach out the payment providers and download
// payments
func GetPayments() {
	payments, err := GetLastMonthsPayments()
	if err != nil {
		log.Errorf("error getting payments: %s", err.Error())
	}

	for _, p := range payments {
		log.Printf("%s, %s, %s", p.Email, p.Amount.Display(), p.Date.String())
	}

	// TODO: add payments to DB
}

// GetLastMonthsPayments fetches payments from paypal
//  to see if members have paid their dues
func GetLastMonthsPayments() ([]Payment, error) {
	startDate := time.Now().AddDate(0, -1, 0).Format(time.RFC3339) // subtract one month
	endDate := time.Now().Format(time.RFC3339)

	p, err := getPaypalPayments(startDate, endDate)
	if err != nil {
		log.Errorf("error getting payments %s", err.Error())
		return p, err
	}

	return mapMemberIDToPayments(p), err
}

// GetLastYearsPayments fetches payments from paypal
//  this is to populate the db with members
func GetLastYearsPayments() ([]Payment, error) {
	startDate := time.Now().AddDate(-1, 0, 0).Format(time.RFC3339) // subtract one year
	endDate := time.Now().Format(time.RFC3339)

	p, err := getPaypalPayments(startDate, endDate)
	if err != nil {
		log.Errorf("error getting payments %s", err.Error())
		return p, err
	}

	return mapMemberIDToPayments(p), err
}

func mapMemberIDToPayments(payments []Payment) []Payment {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}

	for _, p := range payments {
		// if there's no email address, this is not a member payment
		if p.Email == "" {
			continue
		}

		// check that member exists in db
		_, err := db.GetMemberByEmail(p.Email)
		if err != nil {
			log.Errorf("could not find member id from email: %s", err.Error())
			// if member doesn't exist, add them
			_, err := db.AddMember(p.Email, p.Name)
			if err != nil {
				log.Errorf("error adding member to DB: %s", err.Error())
			}
		}
	}

	db.Release()

	return payments
}
