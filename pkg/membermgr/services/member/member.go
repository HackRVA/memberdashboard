package member

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/pkg/membermgr/datastore"
	"github.com/HackRVA/memberserver/pkg/membermgr/integrations"
	"github.com/HackRVA/memberserver/pkg/membermgr/services"
	"github.com/HackRVA/memberserver/pkg/membermgr/services/logger"
	"github.com/HackRVA/memberserver/pkg/slack"

	"github.com/HackRVA/memberserver/pkg/membermgr/models"
)

type member struct {
	model   models.Member
	store   datastore.MemberStore
	service services.Member
}

func NewMemberService(s datastore.MemberStore, m models.Member) member {
	return member{
		model: m,
		store: s,
	}
}

func (m member) PaymentIsBeforeOneMonthAgo(payment models.Payment) bool {
	oneMonthAgo := (time.Hour * 24) * -30
	return payment.Time.Before(time.Now().Add(oneMonthAgo))
}

func (m member) IsActive() bool {
	return m.model.Level == uint8(models.Standard) || m.model.Level == uint8(models.Classic) || m.model.Level == uint8(models.Premium)
}

func (m member) IsCredited() bool {
	return m.model.Level == uint8(models.Credited)
}

func (m member) HasValidSubscriptionID() bool {
	return m.model.SubscriptionID != "none" && len(m.model.SubscriptionID) > 0
}

func (m member) notifyGracePeriod() {
	logger.Infof("[scheduled-job] %s notify about being in grace period", m.model.Name)
	go slack.Send(config.Get().SlackAccessEvents, fmt.Sprintf("%s is in a grace period until their subscription ends", m.model.Name))
}

func (m member) endGracePeriod() {
	logger.Infof("[scheduled-job] %s notify about grace period ending", m.model.Name)
	go slack.Send(config.Get().SlackAccessEvents, fmt.Sprintf("%s grace period has ended. Setting membership level to inactive.", m.model.Name))
}

func (m member) setInactive() {
	logger.Infof("[scheduled-job] %s setting member to inactive", m.model.Name)
	m.store.SetMemberLevel(m.model.ID, models.Inactive)
}

func (m member) UpdateName(name string) {
	if strings.TrimSpace(m.model.Name) != "" {
		return
	}

	if strings.TrimSpace(name) == "" {
		return
	}

	logger.Infof("attempting to update name [%s] from payment provider", name)

	if err := m.service.Update(models.Member{
		ID:   m.model.ID,
		Name: name,
	}); err != nil {
		logger.Error(err)
	}
}

func (m member) UpdateEmail(email string) {
	if strings.TrimSpace(m.model.Email) != "" {
		return
	}

	if strings.TrimSpace(email) == "" {
		return
	}

	logger.Infof("attempting to update email [%s] from payment provider", email)
	if err := m.service.Update(models.Member{
		ID:    m.model.ID,
		Email: email,
	}); err != nil {
		logger.Error(err)
	}
}

func (m member) UpdateInfo(paymentProvider integrations.PaymentProvider) {
	name, email, err := paymentProvider.GetSubscriber(m.model.SubscriptionID)
	if err != nil {
		logger.Error(err)
		return
	}

	if strings.TrimSpace(email) == "" {
		logger.Debugf("did not receive email from payment provider subscription_id: %s", m.model.SubscriptionID)
		return
	}

	if strings.TrimSpace(name) == "" {
		logger.Debugf("did not receive name from payment provider subscription_id: %s", m.model.SubscriptionID)
		return
	}

	logger.Infof("attempting to update member name and email: %s, %s", name, email)

	if err := m.store.UpdateMemberBySubscriptionID(m.model.SubscriptionID, models.Member{
		SubscriptionID: m.model.SubscriptionID,
		Name:           name,
		Email:          email,
	}); err != nil {
		logger.Error(err)
	}
}

func (m member) activeStatusHandler(lastPayment models.Payment) {
	lastPaymentAmount, err := strconv.ParseFloat(lastPayment.Amount, 32)
	if err != nil {
		logger.Error(err)
		return
	}

	if int64(lastPaymentAmount) == models.MemberLevelToAmount[models.Premium] {
		m.store.SetMemberLevel(m.model.ID, models.Premium)
		return
	}
	if int64(lastPaymentAmount) == models.MemberLevelToAmount[models.Classic] {
		m.store.SetMemberLevel(m.model.ID, models.Classic)
		return
	}
	m.store.SetMemberLevel(m.model.ID, models.Standard)
}

func (m member) cancelStatusHandler(lastPayment models.Payment) {
	if m.PaymentIsBeforeOneMonthAgo(lastPayment) {
		if m.IsActive() {
			m.endGracePeriod()
		}
		m.setInactive()

		return
	}
	m.notifyGracePeriod()
}

func (m member) setMemberLevelFromLastPayment(status string, lastPayment models.Payment) {
	logger.Infof("[scheduled-job] setting member level: %s - %s - last payment amount: %s", m.model.Name, status, lastPayment.Amount)

	println(status)
	switch status {
	case models.ActiveStatus:
		m.activeStatusHandler(lastPayment)
		return
	case models.CanceledStatus:
		m.cancelStatusHandler(lastPayment)
		return
	case models.SuspendedStatus:
		m.store.SetMemberLevel(m.model.ID, models.Inactive)
	default:
		return
	}
}

func (m member) CheckStatus(paymentProvider integrations.PaymentProvider) error {
	if m.IsCredited() {
		return nil
	}

	if !m.HasValidSubscriptionID() {
		m.store.SetMemberLevel(m.model.ID, models.Inactive)
		return fmt.Errorf("deactivating member (name: %s email: %s) because no subscriptionID was found", m.model.Name, m.model.Email)
	}

	m.UpdateInfo(paymentProvider)

	status, lastPaymentAmount, lastPaymentTime, err := paymentProvider.GetSubscription(m.model.SubscriptionID)
	if err != nil {
		if !m.IsActive() {
			logger.Debugf("error getting subscription status for (%s, %s). However, member is already inactive. %s", m.model.Email, m.model.Name, err.Error())
			return fmt.Errorf("error getting member's subscription, but the member is already inactive")
		}
		m.store.SetMemberLevel(m.model.ID, models.Inactive)
		return fmt.Errorf("error getting subscription: %s (%s, %s) setting to inactive until status is investigated", err.Error(), m.model.Email, m.model.Name)
	}

	m.setMemberLevelFromLastPayment(status, models.Payment{
		Amount: lastPaymentAmount,
		Time:   lastPaymentTime,
	})
	return nil
}
