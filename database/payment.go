package database

import (
	"context"
	"fmt"

	//"memberserver/mail"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
	log "github.com/sirupsen/logrus"
)

// memberGracePeriod if the member has made a payment in the last 45 days they will remain in a grace period
const memberGracePeriod = 46
const membershipMonth = 31

var paymentDbMethod PaymentDatabaseMethod

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
	rows, err := db.getConn().Query(context.Background(), paymentDbMethod.getPayments())
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

	err := db.getConn().QueryRow(context.Background(), paymentDbMethod.insertPayment(), payment.Date, payment.Amount.AsMajorUnits(), payment.MemberID).Scan(&p.ID, &p.Date, &amount, &p.MemberID)
	if err != nil {
		return fmt.Errorf("conn.Query failed: %v", err)
	}

	p.Amount = *money.New(amount*100, "USD")

	return err
}

// AddPayments adds multiple payments to the database
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

// membersContains checks that a member exists within a list of members
func membersContains(members []Member, m Member) bool {
	for _, v := range members {
		if v.ID == m.ID {
			log.Debugf("%s gets a credit", m.ID)
			return true
		}
	}

	return false
}

// EvaluateMembers loops through all members and evaluates their status
func (db *Database) EvaluateMembers() {
	members := db.GetMembers()

	memberCredits := db.GetMembersWithCredit()
	for _, m := range memberCredits {
		rows, err := db.getConn().Query(context.Background(), paymentDbMethod.updateMembershipLevel(), m.ID, Credited)
		if err != nil {
			log.Errorf("member credit failed: %v", err)
		}
		defer rows.Close()
	}

	for _, m := range members {
		if membersContains(memberCredits, m) {
			continue
		}
		err := db.EvaluateMemberStatus(m.ID)
		if err != nil {
			log.Errorf("error evaluating member's status: %s", err.Error())
		}
	}

	db.clearIdleConnections()
}

// EvaluateMemberStatus look in the db and determine the members' last payment date
//  if it's greater than a certain date revoke their membership
func (db *Database) EvaluateMemberStatus(memberID string) error {
	var daysSincePayment int64
	var amount int64
	var email string

	m, err := db.GetMemberByID(memberID)
	if MemberLevel(m.Level) == Credited {
		return fmt.Errorf("member is credited a membership: %s", err)
	}

	err = db.getConn().QueryRow(context.Background(), paymentDbMethod.checkLastPayment(), memberID).Scan(&daysSincePayment, &amount, &email)
	if err != nil {
		log.Errorf("error finding members payments: %s", err)
	}

	if daysSincePayment > memberGracePeriod { // revoke
		// if this is an active member, then they are just now being revoked.
		// if they are not active, they have already been revoked.
		if MemberLevel(m.Level) == Inactive {
			log.Debugf("Member is already inactive: %s", m.Name)
			return nil
		}
		//TODO: [ML] Move business logic out of database package and user mailer
		//mail.SendRevokedEmail(m.Email, m)

		rows, err := db.getConn().Query(context.Background(), paymentDbMethod.updateMembershipLevel(), memberID, Inactive)
		if err != nil {
			return fmt.Errorf("conn.Query failed: %v", err)
		}
		defer rows.Close()
	} else if daysSincePayment <= memberGracePeriod {
		if daysSincePayment > membershipMonth {

			// currently these would send everyday and everytime the app starts.
			//   it would be better if we could send these only once.

			//c, _ := config.Load()
			// send notification because they are in a grace period
			//TODO: [ML] Move business logic out of database package and user mailer
			// mail.SendGracePeriodMessage(m.Email, m)
			// mail.SendGracePeriodMessageToLeadership(c.AdminEmail, m)
		}

		// a valid member
		// determine their membership level
		mLevel := MemberLevelFromAmount[amount]

		rows, err := db.getConn().Query(context.Background(), paymentDbMethod.updateMembershipLevel(), memberID, mLevel)
		if err != nil {
			return fmt.Errorf("conn.Query failed: %v", err)
		}
		defer rows.Close()
	}

	return nil
}
