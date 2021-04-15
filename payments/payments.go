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
	defer db.Release()

	err = getLastMonthsPayments()
	if err != nil {
		log.Errorf("error getting payments: %s", err.Error())
	}
	log.Debug("done adding payments to db")
}

// getLastMonthsPayments fetches payments from paypal
//  to see if members have paid their dues
func getLastMonthsPayments() error {
	startDate := time.Now().AddDate(0, 0, -15).Format(time.RFC3339) // subtract one month
	endDate := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)     // add a day

	p, err := getPaypalPayments(startDate, endDate)
	if err != nil {
		log.Errorf("error getting payments %s", err.Error())
		return err
	}

	processPayments(p)
	return err
}

func processPayments(payments []database.Payment) {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}
	defer db.Release()

	var membersToAdd []database.Member

	for _, p := range payments {
		if p.Name == "" && p.Email == "" {
			continue
		}

		newMember := database.Member{
			Name:  p.Name,
			Email: p.Email,
		}

		membersToAdd = append(membersToAdd, newMember)
	}

	err = db.AddMembers(membersToAdd)
	if err != nil {
		log.Error(err)
	}

	members := db.GetMembers()

	memberLookup := make(map[string]database.Member)

	for _, m := range members {
		memberLookup[m.Email] = m
	}

	var paymentsWithMemberID []database.Payment
	for _, p := range payments {
		payment := p
		payment.MemberID = memberLookup[p.Email].ID
		paymentsWithMemberID = append(paymentsWithMemberID, payment)
	}

	db.AddPayments(paymentsWithMemberID)
}
