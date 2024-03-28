package member

import (
	"errors"
	"fmt"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/pkg/membermgr/datastore"
	"github.com/HackRVA/memberserver/pkg/membermgr/integrations"
	"github.com/HackRVA/memberserver/pkg/membermgr/services"
	"github.com/HackRVA/memberserver/pkg/slack"

	"github.com/HackRVA/memberserver/pkg/membermgr/models"
)

type memberService struct {
	store           datastore.MemberStore
	resourceManager services.Resource
	paymentProvider integrations.PaymentProvider
	logger          services.Logger
}

func New(store datastore.MemberStore, rm services.Resource, pp integrations.PaymentProvider, logger services.Logger) memberService {
	return memberService{
		store:           store,
		resourceManager: rm,
		paymentProvider: pp,
		logger:          logger,
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
	defer m.CheckStatus(member.SubscriptionID)
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

	if m.SubscriptionID != subscriptionID {
		return m, fmt.Errorf("subscriptionID doesn't match with member: %s, %s", m.Email, m.Name)
	}

	return m, err
}

func (ms memberService) CheckStatus(subscriptionID string) (models.Member, error) {
	var m models.Member

	for _, el := range ms.store.GetMembers() {
		if el.SubscriptionID != subscriptionID {
			continue
		}

		m = el
	}

	if m.SubscriptionID != subscriptionID {
		return m, fmt.Errorf("could not find a member with subscriptionID: %s", subscriptionID)
	}

	mem := member{
		model:   m,
		store:   ms.store,
		service: ms,
	}

	return m, mem.CheckStatus(ms.paymentProvider)
}

func (m memberService) GetTiers() []models.Tier {
	return m.store.GetTiers()
}

func (m memberService) FindNonMembersOnSlack() []string {
	var nonMembers []string
	users, err := slack.GetUsers(config.Get().SlackToken)
	if err != nil {
		m.logger.Errorf("error fetching slack users: %s", err)
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
