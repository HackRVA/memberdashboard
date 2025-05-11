package member

import (
	"errors"
	"fmt"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/datastore"
	"github.com/HackRVA/memberserver/integrations"
	"github.com/HackRVA/memberserver/pkg/slack"
	"github.com/HackRVA/memberserver/services"
	"github.com/HackRVA/memberserver/services/logger"
	"github.com/sirupsen/logrus"

	"github.com/HackRVA/memberserver/models"
)

type memberService struct {
	store           datastore.MemberStore
	resourceManager services.ResourceUpdater
	paymentProvider integrations.PaymentProvider
}

func New(store datastore.MemberStore, rm services.ResourceUpdater, pp integrations.PaymentProvider) memberService {
	return memberService{
		store:           store,
		resourceManager: rm,
		paymentProvider: pp,
	}
}

func (m memberService) Add(newMember models.Member) (models.Member, error) {
	createdMember, err := m.store.AddNewMember(newMember)
	if err != nil {
		logrus.Error(err)
	}

	// assignRFID needs to run after the member has been added to the DB
	memberWithRFID, err := m.AssignRFID(createdMember.Email, createdMember.RFID) 
	if err != nil {
		logrus.Error(err)
	}

	createdMember.RFID = memberWithRFID.RFID

	return createdMember, nil
}

func (m memberService) GetMembersPaginated(limit int, count int, active bool) []models.Member {
	members, err := m.store.GetMembersPaginated(limit, count, active)
	if err != nil {
		logrus.Error(err)
	}
	return members
}

func (m memberService) GetMemberCount(isActive bool) (int, error) {
	return m.store.GetMemberCount(isActive)
}

func (m memberService) Get() []models.Member {
	return m.store.GetMembers()
}

func (m memberService) GetByEmail(email string) (models.Member, error) {
	return m.store.GetMemberByEmail(email)
}

func (m memberService) Update(member models.Member) error {
	if _, err := m.CheckStatus(member.SubscriptionID); err != nil {
		logrus.Error(err)
	}
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

func (ms memberService) GetMemberBySubscriptionID(subscriptionID string) (models.Member, error) {
	_, email, err := ms.paymentProvider.GetSubscriber(subscriptionID)
	if err != nil {
		return models.Member{}, err
	}

	m, err := ms.GetByEmail(email)
	if err != nil {
		return models.Member{}, err
	}

	if m.SubscriptionID != subscriptionID {
		logrus.Errorf("subscriptionID doesn't match with member: %s, %s", m.Email, m.Name)
		return m, fmt.Errorf("subscriptionID doesn't match with member: %s, %s", m.Email, m.Name)
	}

	return m, nil
}

func (ms memberService) CheckStatus(subscriptionID string) (models.Member, error) {
	var m models.Member

	if subscriptionID == "none" {
		logrus.Error("tried to lookup subscriptionID that was 'none'")
		return m, errors.New("tried to lookup subscriptionID that was 'none'")
	}

	for _, el := range ms.store.GetMembers() {
		if el.SubscriptionID != subscriptionID {
			continue
		}

		m = el
	}

	if m.SubscriptionID != subscriptionID {
		return m, fmt.Errorf("could not find a member with subscriptionID: %s", subscriptionID)
	}

	statusChecker := NewStatusChecker(m, ms.store, ms.paymentProvider)

	return m, statusChecker.CheckStatus()
}

func (m memberService) GetTiers() []models.Tier {
	return m.store.GetTiers()
}

func (m memberService) FindNonMembersOnSlack() []string {
	var nonMembers []string
	users, err := slack.GetUsers(config.Get().SlackToken)
	if err != nil {
		logger.Errorf("error fetching slack users: %s", err)
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

func (ms memberService) SetLevel(memberID string, level models.MemberLevel) error {
	return ms.store.SetMemberLevel(memberID, level)
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

func (ms memberService) GetActiveMembersWithoutSubscription() []models.Member {
	return ms.store.GetActiveMembersWithoutSubscription()
}
