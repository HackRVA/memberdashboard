package payments

import (
	"fmt"
	"memberserver/database"
	"sync"
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
	startDate := time.Now().AddDate(0, -1, 0).Format(time.RFC3339) // subtract one month
	endDate := time.Now().Format(time.RFC3339)

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

	var wg sync.WaitGroup
	var awg sync.WaitGroup
	paymentChan := make(chan database.Payment, len(payments)+1)

	for _, p := range payments {
		wg.Add(1)
		go processPayment(p, paymentChan, db, &wg)
	}
	wg.Wait()
	log.Debug("done processing payments")

	for range payments {
		select {
		case pay := <-paymentChan:
			awg.Add(1)
			go addPaymentToDB(pay, db, &awg)
		default:
			fmt.Println("payment wasn't received in channel")
		}
	}
	awg.Wait()
	close(paymentChan)
}

// processPayment will attribute a payment to a member.
func processPayment(p database.Payment, paymentChan chan database.Payment, db *database.Database, wg *sync.WaitGroup) {
	defer wg.Done()

	// if there's no email address, this is not a member payment
	// this could possibly be outgoing transactions
	if p.Email == "" {
		log.Debugf("transaction without email: %s %s %s %s\n", p.Amount.Display(), p.Name, p.Email, p.Date.String())
		return
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
		paymentChan <- p
		return
	}

	p.MemberID = m.ID
	paymentChan <- p
}

func addPaymentToDB(p database.Payment, db *database.Database, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Debugf("adding payment to db: %s %s %s %s\n", p.Amount.Display(), p.Name, p.Email, p.Date.String())

	err := db.AddPayment(p)
	if err != nil {
		log.Errorf("error adding payment to db: %s, %s, %s, %s: %s", p.Email, p.Amount.Display(), p.Date.String(), p.MemberID, err.Error())
	}
}
