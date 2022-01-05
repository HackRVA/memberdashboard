package payments

import (
	"memberserver/api/models"
	"memberserver/datastore"
	"time"

	log "github.com/sirupsen/logrus"
)

type PaymentProvider struct {
	db          datastore.DataStore
	accessToken string
}

func Setup(database datastore.DataStore) PaymentProvider {
	return PaymentProvider{
		db: database,
	}
}

// GetPayments reach out the payment providers and download
// payments
func (p PaymentProvider) GetPayments() {
	err := p.getLastMonthsPayments()
	if err != nil {
		log.Errorf("error getting payments: %s", err.Error())
	}
	log.Debug("done adding payments to db")
}

// getLastMonthsPayments fetches payments from paypal
//  to see if members have paid their dues
func (p PaymentProvider) getLastMonthsPayments() error {
	startDate := time.Now().AddDate(0, 0, -15).Format(time.RFC3339) // subtract one month
	endDate := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)     // add a day

	payment, err := p.getPaypalPayments(startDate, endDate)
	if err != nil {
		log.Errorf("error getting payments %s", err.Error())
		return err
	}

	p.processPayments(payment)
	return err
}

func (p PaymentProvider) processPayments(payments []models.Payment) {
	var membersToAdd []models.Member

	for _, p := range payments {
		if p.Name == "" && p.Email == "" {
			continue
		}

		newMember := models.Member{
			Name:  p.Name,
			Email: p.Email,
		}

		membersToAdd = append(membersToAdd, newMember)
	}

	err := p.db.AddMembers(membersToAdd)
	if err != nil {
		log.Error(err)
	}

	members := p.db.GetMembers()

	memberLookup := make(map[string]models.Member)

	for _, m := range members {
		memberLookup[m.Email] = m
	}

	var paymentsWithMemberID []models.Payment
	for _, p := range payments {
		payment := p
		payment.MemberID = memberLookup[p.Email].ID
		paymentsWithMemberID = append(paymentsWithMemberID, payment)
	}

	p.db.AddPayments(paymentsWithMemberID)
}
