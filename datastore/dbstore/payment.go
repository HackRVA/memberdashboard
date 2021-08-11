package dbstore

import (
	"context"
	"fmt"
	"memberserver/api/models"
	"strings"

	"github.com/Rhymond/go-money"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

var paymentDbMethod PaymentDatabaseMethod

// GetPayments - get list of payments that we have in the db
func (db *DatabaseStore) GetPayments() ([]models.Payment, error) {
	rows, err := db.getConn().Query(context.Background(), paymentDbMethod.getPayments())
	if err != nil {
		log.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	var payments []models.Payment

	for rows.Next() {
		var p models.Payment
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
func (db *DatabaseStore) AddPayment(payment models.Payment) error {
	var p models.Payment
	var amount int64

	err := db.getConn().QueryRow(context.Background(), paymentDbMethod.insertPayment(), payment.Date, payment.Amount.AsMajorUnits(), payment.MemberID).Scan(&p.ID, &p.Date, &amount, &p.MemberID)
	if err != nil {
		return fmt.Errorf("conn.Query failed: %v", err)
	}

	p.Amount = *money.New(amount*100, "USD")

	return err
}

// AddPayments adds multiple payments to the database
func (db *DatabaseStore) AddPayments(payments []models.Payment) error {
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

	_, err := db.getConn().Exec(context.Background(), sqlStr+str+" ON CONFLICT DO NOTHING;")
	if err != nil {
		return fmt.Errorf("conn.Exec failed: %v", err)
	}

	return err
}

// SetMemberLevel sets a member's membership tier
func (db *DatabaseStore) SetMemberLevel(memberId string, level models.MemberLevel) error {
	rows, err := db.getConn().Query(context.Background(), paymentDbMethod.updateMembershipLevel(), memberId, level)
	if err != nil {
		log.Errorf("Set member level failed: %v", err)
		return err
	}
	defer rows.Close()
	return nil
}

// ApplyMemberCredits updates members tiers for all members with credit to Credited
func (db *DatabaseStore) ApplyMemberCredits() {
	//	Member credits are currently managed by DB commands.  #102 will address this.
	memberCredits := db.GetMembersWithCredit()
	for _, m := range memberCredits {
		err := db.SetMemberLevel(m.ID, models.Credited)
		if err != nil {
			log.Errorf("member credit failed: %v", err)
		}
	}
}

// UpdateMemberTiers updates member tiers based on the most recent payment amount
func (db *DatabaseStore) UpdateMemberTiers() {
	db.getConn().Exec(context.Background(), paymentDbMethod.updateMemberTiers())
}

// GetPastDueAccounts retrieves all active members without a payment in the last month
func (db *DatabaseStore) GetPastDueAccounts() []models.PastDueAccount {
	var pastDueAccounts []models.PastDueAccount
	rows, err := db.getConn().Query(context.Background(), paymentDbMethod.pastDuePayments())

	if err == pgx.ErrNoRows {
		return pastDueAccounts
	}

	if err != nil {
		log.Errorf("conn.Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.PastDueAccount
		err = rows.Scan(&p.MemberId, &p.Name, &p.Email, &p.LastPaymentDate, &p.DaysSinceLastPayment)
		if err != nil {
			log.Errorf("error scanning row: %s", err)
		}
		pastDueAccounts = append(pastDueAccounts, p)
	}

	return pastDueAccounts
}
