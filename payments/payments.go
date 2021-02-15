package payments

import (
	"memberserver/database"
	"time"

	log "github.com/sirupsen/logrus"
)

type paypalAuth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// GetPayments reach out the payment providers and download
// payments
func GetPayments() {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}

	payments, err := getLastMonthsPayments()
	if err != nil {
		log.Errorf("error getting payments: %s", err.Error())
	}

	for _, p := range payments {
		_, err := db.AddPayment(p)
		if err != nil {
			log.Errorf("error adding payment to db: %s, %s, %s, %s: %s", p.Email, p.Amount.Display(), p.Date.String(), p.MemberID, err.Error())
		}
	}

	db.Release()
}

// getLastMonthsPayments fetches payments from paypal
//  to see if members have paid their dues
func getLastMonthsPayments() ([]database.Payment, error) {
	startDate := time.Now().AddDate(0, -1, 0).Format(time.RFC3339) // subtract one month
	endDate := time.Now().Format(time.RFC3339)

	p, err := getPaypalPayments(startDate, endDate)
	if err != nil {
		log.Errorf("error getting payments %s", err.Error())
		return p, err
	}

	return mapMemberIDToPayments(p), err
}

func mapMemberIDToPayments(payments []database.Payment) []database.Payment {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}

	var paymentsWithIDs []database.Payment
	for _, p := range payments {
		// if there's no email address, this is not a member payment
		if p.Email == "" {
			continue
		}

		// check that member exists in db
		m, err := db.GetMemberByEmail(p.Email)
		if err != nil {
			// if member doesn't exist, add them
			am, err := db.AddMember(p.Email, p.Name)
			if err != nil {
				log.Errorf("error adding member to DB: %s", err.Error())
			}
            db.AddUserToDefaultResources(p.Email)
			p.MemberID = am.ID
		} else {
			p.MemberID = m.ID
		}
		paymentsWithIDs = append(paymentsWithIDs, p)
	}

	db.Release()

	return paymentsWithIDs
}
