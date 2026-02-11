package member

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/datastore"
	"github.com/HackRVA/memberserver/integrations"
	"github.com/HackRVA/memberserver/models"
	"github.com/HackRVA/memberserver/pkg/slack"
	"github.com/HackRVA/memberserver/services/logger"
)

type StatusChecker struct {
	member          models.Member
	store           datastore.MemberStore
	paymentProvider integrations.PaymentProvider
}

func NewStatusChecker(member models.Member, store datastore.MemberStore, paymentProvider integrations.PaymentProvider) *StatusChecker {
	return &StatusChecker{
		member:          member,
		store:           store,
		paymentProvider: paymentProvider,
	}
}

func (s *StatusChecker) PaymentIsBeforeOneMonthAgo(payment models.Payment) bool {
	oneMonthAgo := (time.Hour * 24) * -30
	return payment.Time.Before(time.Now().Add(oneMonthAgo))
}

func (s *StatusChecker) IsActive() bool {
	return s.member.Level == uint8(models.Standard) || s.member.Level == uint8(models.Classic) || s.member.Level == uint8(models.Premium)
}

func (s *StatusChecker) IsCredited() bool {
	return s.member.Level == uint8(models.Credited)
}

func (s *StatusChecker) HasValidSubscriptionID() bool {
	return s.member.SubscriptionID != "none" && len(s.member.SubscriptionID) > 0
}

func (s *StatusChecker) notifyGracePeriod(lastPayment models.Payment) {
	logger.Infof("[scheduled-job] %s notify about being in grace period. Last payment date: %s", s.member.Name, lastPayment.Time.Format("2006-01-02"))
	go slack.Send(config.Get().SlackAccessEvents, fmt.Sprintf("%s is in a grace period until their subscription ends. \n\tLast payment date: %s", s.member.Name, lastPayment.Time.Format("2006-01-02")))
}

func (s *StatusChecker) endGracePeriod() {
	logger.Infof("[scheduled-job] %s notify about grace period ending", s.member.Name)
	go slack.Send(config.Get().SlackAccessEvents, fmt.Sprintf("%s grace period has ended. Setting membership level to inactive.", s.member.Name))
}

func (s *StatusChecker) setInactive() {
	logger.Infof("[scheduled-job] %s setting member to inactive", s.member.Name)
	if err := s.store.SetMemberLevel(s.member.ID, models.Inactive); err != nil {
		logger.Errorf("error setting member level %s", err)
	}
}

func (s *StatusChecker) UpdateName(name string) {
	if strings.TrimSpace(s.member.Name) != "" {
		return
	}

	if strings.TrimSpace(name) == "" {
		return
	}

	logger.Infof("attempting to update name [%s] from payment provider", name)

	if err := s.store.UpdateMember(models.Member{
		ID:   s.member.ID,
		Name: name,
	}); err != nil {
		logger.Error(err)
	}
}

func (s *StatusChecker) UpdateEmail(email string) {
	if strings.TrimSpace(s.member.Email) != "" {
		return
	}

	if strings.TrimSpace(email) == "" {
		return
	}

	logger.Infof("attempting to update email [%s] from payment provider", email)
	if err := s.store.UpdateMember(models.Member{
		ID:    s.member.ID,
		Email: email,
	}); err != nil {
		logger.Error(err)
	}
}

func (s *StatusChecker) UpdateInfo() {
	name, email, err := s.paymentProvider.GetSubscriber(s.member.SubscriptionID)
	if err != nil {
		logger.Error(err)
		return
	}

	if strings.TrimSpace(email) == "" {
		logger.Debugf("did not receive email from payment provider subscription_id: %s", s.member.SubscriptionID)
		return
	}

	if strings.TrimSpace(name) == "" {
		logger.Debugf("did not receive name from payment provider subscription_id: %s", s.member.SubscriptionID)
		return
	}

	logger.Infof("attempting to update member name and email: %s, %s", name, email)

	if err := s.store.UpdateMemberBySubscriptionID(s.member.SubscriptionID, models.Member{
		SubscriptionID: s.member.SubscriptionID,
		Name:           name,
		Email:          email,
	}); err != nil {
		logger.Error(err)
	}
}

func (s *StatusChecker) activeStatusHandler(lastPayment models.Payment) {
	lastPaymentAmount, err := strconv.ParseFloat(lastPayment.Amount, 32)
	if err != nil {
		logger.Error(err)
		return
	}

	if int64(lastPaymentAmount) == models.MemberLevelToAmount[models.Premium] {
		if err := s.store.SetMemberLevel(s.member.ID, models.Premium); err != nil {
			logger.Error(err)
		}
		return
	}
	if int64(lastPaymentAmount) == models.MemberLevelToAmount[models.Classic] {
		if err := s.store.SetMemberLevel(s.member.ID, models.Classic); err != nil {
			logger.Error(err)
		}
		return
	}
	if err := s.store.SetMemberLevel(s.member.ID, models.Standard); err != nil {
		logger.Error(err)
	}
}

func (s *StatusChecker) cancelStatusHandler(lastPayment models.Payment) {
	if s.PaymentIsBeforeOneMonthAgo(lastPayment) {
		if s.IsActive() {
			s.endGracePeriod()
		}
		s.setInactive()

		return
	}
	s.notifyGracePeriod(lastPayment)
}

func (s *StatusChecker) setMemberLevelFromLastPayment(status string, lastPayment models.Payment) {
	logger.Infof("[scheduled-job] setting member level: %s - %s - last payment amount: %s", s.member.Name, status, lastPayment.Amount)

	switch status {
	case models.ActiveStatus:
		s.activeStatusHandler(lastPayment)
		return
	case models.CanceledStatus:
		s.cancelStatusHandler(lastPayment)
		return
	case models.SuspendedStatus:
		if err := s.store.SetMemberLevel(s.member.ID, models.Inactive); err != nil {
			logger.Error(err)
		}
	default:
		return
	}
}

// CheckStatus looks at the embedded Member on `StatusChecker`
// it will verify that they have a valid subscriptionID.
// it will determine the appropriate member status based on
// last payment date and subscription status.
func (s *StatusChecker) CheckStatus() error {
	if s.IsCredited() {
		return nil
	}

	if !s.HasValidSubscriptionID() {
		if err := s.store.SetMemberLevel(s.member.ID, models.Inactive); err != nil {
			logger.Error(err)
		}
		return fmt.Errorf("deactivating member (name: %s email: %s) because no subscriptionID was found", s.member.Name, s.member.Email)
	}

	// s.UpdateInfo()

	status, lastPaymentAmount, lastPaymentTime, err := s.paymentProvider.GetSubscription(s.member.SubscriptionID)
	if err != nil {
		if !s.IsActive() {
			logger.Debugf("error getting subscription status for (%s, %s). However, member is already inactive. %s", s.member.Email, s.member.Name, err.Error())
			return fmt.Errorf("error getting member's subscription, but the member is already inactive")
		}
		if err := s.store.SetMemberLevel(s.member.ID, models.Inactive); err != nil {
			logger.Error(err)
		}
		return fmt.Errorf("error getting subscription: %s (%s, %s) setting to inactive until status is investigated", err.Error(), s.member.Email, s.member.Name)
	}

	s.setMemberLevelFromLastPayment(status, models.Payment{
		Amount: lastPaymentAmount,
		Time:   lastPaymentTime,
	})
	return nil
}
