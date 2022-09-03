package member

import (
	"errors"
	"fmt"
	"memberserver/internal/datastore"
	"memberserver/internal/integrations"
	"memberserver/internal/models"
	"memberserver/internal/services/resourcemanager"
	"memberserver/pkg/slack"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type memberService struct {
	store           datastore.MemberStore
	resourceManager resourcemanager.ResourceManager
	paymentProvider integrations.PaymentProvider
}

type MemberService interface {
	Add(models.Member) (models.Member, error)
	Get() []models.Member
	GetMembersWithLimit(limit int, offset int, active bool) []models.Member
	GetByEmail(email string) (models.Member, error)
	Update(models.Member) error
	AssignRFID(email string, rfid string) (models.Member, error)
	GetTiers() []models.Tier
	FindNonMembersOnSlack() []string
	CancelStatusHandler(member models.Member, lastPayment models.Payment)
	ActiveStatusHandler(member models.Member, lastPayment models.Payment)
	GetMemberFromSubscription(subscriptionID string) (models.Member, error)
}

func NewMemberService(store datastore.MemberStore, rm resourcemanager.ResourceManager, pp integrations.PaymentProvider) memberService {
	return memberService{
		store:           store,
		resourceManager: rm,
		paymentProvider: pp,
	}
}

func (m memberService) Add(newMember models.Member) (models.Member, error) {
	// assignRFID needs to run after the member has been added to the DB
	defer m.AssignRFID(newMember.Email, newMember.RFID)
	return m.store.AddNewMember(newMember)
}

func (m memberService) GetMembersWithLimit(limit int, count int, active bool) []models.Member {
	return m.store.GetMembersWithLimit(limit, count, active)
}
func (m memberService) Get() []models.Member {
	return m.store.GetMembers()
}

func (m memberService) GetByEmail(email string) (models.Member, error) {
	return m.store.GetMemberByEmail(email)
}

func (m memberService) Update(member models.Member) error {
	return m.store.UpdateMember(member)
}

func (m memberService) AssignRFID(email string, rfid string) (models.Member, error) {
	if len(rfid) == 0 {
		return models.Member{}, errors.New("not a valid rfid")
	}

	// we need to push to resources after we add rfid to DB
	defer m.resourceManager.PushOne(models.Member{Email: email})
	return m.store.AssignRFID(email, rfid)
}

func (m memberService) GetTiers() []models.Tier {
	return m.store.GetTiers()
}

func (m memberService) FindNonMembersOnSlack() []string {
	var nonMembers []string

	users, err := slack.GetUsers()
	if err != nil {
		logrus.Errorf("error fetching slack users: %s", err)
	}

	members := m.Get()
	memberMap := make(map[string]models.Member)

	for _, m := range members {
		memberMap[m.Email] = m
	}

	for _, u := range users {
		if u.IsBot {
			continue
		}

		if u.Deleted {
			continue
		}

		_, ok := memberMap[u.Profile.Email]
		if !ok {
			nonMembers = append(nonMembers, u.RealName+", "+u.Profile.Email)
		}
	}
	return nonMembers
}

func (ms memberService) CancelStatusHandler(m models.Member, lastPayment models.Payment) {
	member := member{
		store: ms.store,
		model: m,
	}

	if member.paymentIsBeforeOneMonthAgo(lastPayment) {
		if member.isActive() {
			member.endGracePeriod()
		}
		member.setInactive()

		return
	}
	member.notifyGracePeriod()
}

func (ms memberService) ActiveStatusHandler(m models.Member, lastPayment models.Payment) {
	member := member{
		store: ms.store,
		model: m,
	}
	lastPaymentAmount, err := strconv.ParseFloat(lastPayment.Amount, 32)
	if err != nil {
		logrus.Error(err)
	}

	member.setMemberLevel(lastPaymentAmount)
}

func (ms memberService) GetMemberFromSubscription(subscriptionID string) (models.Member, error) {
	name, email, err := ms.paymentProvider.GetSubscriber(subscriptionID)
	if err != nil {
		return models.Member{}, err
	}
	return models.Member{
		Email:          email,
		Name:           name,
		SubscriptionID: subscriptionID,
	}, nil
}

type member struct {
	model models.Member
	store datastore.MemberStore
}

func (m member) paymentIsBeforeOneMonthAgo(payment models.Payment) bool {
	oneMonthAgo := (time.Hour * 24) * -30
	return payment.Time.Before(time.Now().Add(oneMonthAgo))
}

func (m member) isActive() bool {
	return m.model.Level == uint8(models.Standard) || m.model.Level == uint8(models.Classic) || m.model.Level == uint8(models.Premium)
}

func (m member) notifyGracePeriod() {
	logrus.Infof("[scheduled-job] %s notify about being in grace period", m.model.Name)
	go slack.Send(fmt.Sprintf("%s is in a grace period until their subscription ends", m.model.Name))
}

func (m member) endGracePeriod() {
	logrus.Infof("[scheduled-job] %s notify about grace period ending", m.model.Name)
	go slack.Send(fmt.Sprintf("%s grace period has ended. Setting membership level to inactive.", m.model.Name))
}

func (m member) setInactive() {
	logrus.Infof("[scheduled-job] %s setting member to inactive", m.model.Name)
	m.store.SetMemberLevel(m.model.ID, models.Inactive)
}

func (m member) setMemberLevel(lastPaymentAmount float64) {
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
