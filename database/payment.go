package database

import (
	"context"
	"fmt"
	"memberserver/config"
	"memberserver/mail"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
	log "github.com/sirupsen/logrus"
)

// memberGracePeriod if the member has made a payment in the last 45 days they will remain in a grace period
const memberGracePeriod = 46
const membershipMonth = 31

const getPaymentsQuery = `
SELECT id, date, amount
FROM membership.payments
ORDER BY date;`

const insertPaymentQuery = `
INSERT INTO membership.payments(
date, amount, member_id)
VALUES ($1, $2, $3)
RETURNING *;`

// checkRecentPayment - if the member doesn't have a recent payment,
//    we will revoke their membership
const checkLastPaymentQuery = `
SELECT current_date - date as last_payment, amount, email
FROM membership.payments
LEFT JOIN membership.members
ON membership.payments.member_id = membership.members.id
WHERE member_id = $1
ORDER BY date DESC
limit $2;`

const countPaymentsOfMemberSinceQuery = `
SELECT COUNT(*) as num_payments
FROM membership.payments
LEFT JOIN membership.members
ON membership.payments.member_id = membership.members.id
WHERE member_id = $1
AND date >= current_date - $2;`

const updateMembershipLevelQuery = `
UPDATE membership.members
SET member_tier_id=$2
WHERE id=$1
RETURNING *;`

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
	Date     time.Time
	Amount   money.Money
	Provider PaymentProvider
	MemberID string
	Email    string
	Name     string
}

// GetPayments - get list of payments that we have in the db
func (db *Database) GetPayments() ([]Payment, error) {
	rows, err := db.getConn().Query(context.Background(), getPaymentsQuery)
	if err != nil {
		log.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	var payments []Payment

	for rows.Next() {
		var p Payment
		var amount int64
		err = rows.Scan(&p.ID, &p.Date, &amount)
		if err != nil {
			log.Errorf("error scanning row: %s", err)
		}

		p.Amount = *money.New(amount*100, "USD")

		payments = append(payments, p)
	}

	return payments, nil
}

// AddPayment adds a member to the database
func (db *Database) AddPayment(payment Payment) error {
	var p Payment
	var amount int64

	err := db.getConn().QueryRow(context.Background(), insertPaymentQuery, payment.Date, payment.Amount.AsMajorUnits(), payment.MemberID).Scan(&p.ID, &p.Date, &amount, &p.MemberID)
	if err != nil {
		return fmt.Errorf("conn.Query failed: %v", err)
	}

	p.Amount = *money.New(amount*100, "USD")

	return err
}

// AddPayment adds multiple payments to the database
func (db *Database) AddPayments(payments []Payment) error {
	var valStr []string

	sqlStr := `INSERT INTO membership.payments(
date, amount, member_id)
VALUES `

	for _, p := range payments {
		if p.MemberID == "" {
			continue
		}
		valStr = append(valStr, fmt.Sprintf("('%s', %d, '%s')", p.Date.Format("2006-01-02"), p.Amount.Amount()/100, p.MemberID))
	}

	str := strings.Join(valStr, ",")

	_, err := db.getConn().Query(context.Background(), sqlStr+str+" ON CONFLICT DO NOTHING;")
	if err != nil {
		return fmt.Errorf("conn.Query failed: %v", err)
	}

	return err
}

// EvaluateMemberStatus look in the db and determine the members' last payment date
//  if it's greater than a certain date revoke their membership
func (db *Database) EvaluateMemberStatus(memberID string) error {
	var daysSincePayment int64
	var amount int64
	var email string

	// TODO: see if they have multiple memberships
	numMemberships := 1

	err := db.getConn().QueryRow(context.Background(), checkLastPaymentQuery, memberID, numMemberships).Scan(&daysSincePayment, &amount, &email)
	if err != nil {
		return fmt.Errorf("conn.Query failed: %v", err)
	}

	// log.Debugf("days since payment: %d payment amount: %d", daysSincePayment, amount)

	if daysSincePayment > memberGracePeriod { // revoke
		// sendRevokedEmail(email)
		rows, err := db.getConn().Query(context.Background(), updateMembershipLevelQuery, memberID, Inactive)
		if err != nil {
			return fmt.Errorf("conn.Query failed: %v", err)
		}
		defer rows.Close()
	} else if daysSincePayment <= memberGracePeriod {
		if daysSincePayment > membershipMonth {
			// send notification because they are in a grace period
			// sendGracePeriodMessage(email)

			sendGracePeriodMessageToLeadership(email)
		}

		// a valid member
		// determine their membership level
		mLevel := MemberLevelFromAmount[amount]

		rows, err := db.getConn().Query(context.Background(), updateMembershipLevelQuery, memberID, mLevel)
		if err != nil {
			return fmt.Errorf("conn.Query failed: %v", err)
		}
		defer rows.Close()
	}

	return nil
}

func sendGracePeriodMessageToLeadership(address string) {
	conf, _ := config.Load()
	if !conf.EnableInfoEmails {
		return
	}

	mp, err := mail.Setup()
	if err != nil {
		log.Errorf("error setting up mailprovider when attempting to send email notification")
	}

	mp.SendSMTP("info@hackrva.org", address+": hackrva grace period", address+" membership is entering a grace period.")
}

func sendGracePeriodMessage(address string) {
	// mp, err := mail.Setup()
	// if err != nil {
	// 	log.Errorf("error setting up mailprovider when attempting to send email notification")
	// }

	// mp.SendSMTP(address, "hackrva grace period", "you're membership is entering a grace period.  Please try to pay your hackrva membership dues soon.")
}

func sendRevokedEmail(address string) {
	// mp, err := mail.Setup()
	// if err != nil {
	// 	log.Errorf("error setting up mailprovider when attempting to send email notification")
	// }

	// mp.SendSMTP(address, "hackrva membership revoked", "Unfortunately, hackrva hasn't received your membership dues.  Your membership has been revoked until a payment is received.  Please reach out to us if you have any concerns.")
}
