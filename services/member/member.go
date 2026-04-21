package member

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

func (m memberService) Add(ctx context.Context, newMember models.Member) (models.Member, error) {
	createdMember, err := m.store.AddNewMember(ctx, newMember.EnsureUpperSubscriptionID())
	if err != nil {
		logrus.Error(err)
	}

	// assignRFID needs to run after the member has been added to the DB
	memberWithRFID, err := m.AssignRFID(ctx, createdMember.Email, createdMember.RFID)
	if err != nil {
		logrus.Error(err)
	}

	createdMember.RFID = memberWithRFID.RFID

	return createdMember, nil
}

func (m memberService) GetMembersPaginated(ctx context.Context, limit int, count int, active bool) []models.Member {
	members, err := m.store.GetMembersPaginated(ctx, limit, count, active)
	if err != nil {
		logrus.Error(err)
	}
	return members
}

func (m memberService) GetMemberCount(ctx context.Context, isActive bool) (int, error) {
	return m.store.GetMemberCount(ctx, isActive)
}

func (m memberService) Get(ctx context.Context) []models.Member {
	return m.store.GetMembers(ctx)
}

func (m memberService) GetByEmail(ctx context.Context, email string) (models.Member, error) {
	return m.store.GetMemberByEmail(ctx, email)
}

func (m memberService) Update(ctx context.Context, member models.Member) error {
	if _, err := m.CheckStatus(ctx, member.SubscriptionID); err != nil {
		logrus.Error(err)
	}
	return m.store.UpdateMember(ctx, member)
}

func (m memberService) UpdateMemberByID(ctx context.Context, memberID string, update models.Member) error {
	return m.store.UpdateMemberByID(ctx, memberID, update)
}

func (m memberService) AssignRFID(ctx context.Context, email string, rfid string) (models.Member, error) {
	if len(rfid) == 0 {
		return models.Member{}, errors.New("not a valid rfid")
	}

	m.resourceManager.PushOne(models.Member{Email: email})
	return m.store.AssignRFID(ctx, email, rfid)
}

func (ms memberService) GetMemberBySubscriptionID(ctx context.Context, subscriptionID string) (models.Member, error) {
	_, email, err := ms.paymentProvider.GetSubscriber(subscriptionID)
	if err != nil {
		return models.Member{}, err
	}

	m, err := ms.GetByEmail(ctx, email)
	if err != nil {
		return models.Member{}, err
	}

	if !strings.EqualFold(m.SubscriptionID, subscriptionID) {
		logrus.Errorf("subscriptionID doesn't match with member: %s, %s", m.Email, m.Name)
		return m, fmt.Errorf("subscriptionID doesn't match with member: %s, %s", m.Email, m.Name)
	}

	return m, nil
}

func (ms memberService) CheckStatus(ctx context.Context, subscriptionID string) (models.Member, error) {
	var m models.Member

	if subscriptionID == "none" {
		logrus.Error("tried to lookup subscriptionID that was 'none'")
		return m, errors.New("tried to lookup subscriptionID that was 'none'")
	}

	for _, el := range ms.store.GetMembers(ctx) {
		if !strings.EqualFold(el.SubscriptionID, subscriptionID) {
			continue
		}

		m = el
	}

	if !strings.EqualFold(m.SubscriptionID, subscriptionID) {
		return m, fmt.Errorf("could not find a member with subscriptionID: %s", subscriptionID)
	}

	statusChecker := NewStatusChecker(m, ms.store, ms.paymentProvider)

	return m, statusChecker.CheckStatus(ctx)
}

func (m memberService) GetTiers(ctx context.Context) []models.Tier {
	return m.store.GetTiers(ctx)
}

func (m memberService) FindNonMembersOnSlack(ctx context.Context) []string {
	var nonMembers []string
	users, err := slack.GetUsers(config.Get().SlackToken)
	if err != nil {
		logger.Errorf("error fetching slack users: %s", err)
	}

	members := m.Get(ctx)
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

func (ms memberService) SetLevel(ctx context.Context, memberID string, level models.MemberLevel) error {
	return ms.store.SetMemberLevel(ctx, memberID, level)
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

func (ms memberService) GetActiveMembersWithoutSubscription(ctx context.Context) []models.Member {
	return ms.store.GetActiveMembersWithoutSubscription(ctx)
}
