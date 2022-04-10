package member

import (
	"errors"
	"memberserver/internal/datastore"
	"memberserver/internal/models"
	"memberserver/internal/services/resourcemanager"
	"memberserver/pkg/slack"

	"github.com/sirupsen/logrus"
)

type Member struct {
	store           datastore.MemberStore
	resourceManager resourcemanager.ResourceManager
}

type MemberService interface {
	Add(models.Member) (models.Member, error)
	Get() []models.Member
	GetByEmail(email string) (models.Member, error)
	Update(models.Member) error
	AssignRFID(email string, rfid string) (models.Member, error)
	GetTiers() []models.Tier
	FindNonMembersOnSlack() []string
}

func NewMemberService(store datastore.MemberStore, rm resourcemanager.ResourceManager) Member {
	return Member{
		store:           store,
		resourceManager: rm,
	}
}

func (m Member) Add(newMember models.Member) (models.Member, error) {
	m.AssignRFID(newMember.Email, newMember.RFID)
	return m.store.AddNewMember(newMember)
}

func (m Member) Get() []models.Member {
	return m.store.GetMembers()
}

func (m Member) GetByEmail(email string) (models.Member, error) {
	return m.store.GetMemberByEmail(email)
}

func (m Member) Update(member models.Member) error {
	return m.store.UpdateMember(member)
}

func (m Member) AssignRFID(email string, rfid string) (models.Member, error) {
	if len(rfid) == 0 {
		return models.Member{}, errors.New("not a valid rfid")
	}

	m.resourceManager.RemoveOne(models.Member{Email: email})
	go m.resourceManager.PushOne(models.Member{Email: email})
	return m.store.AssignRFID(email, rfid)
}

func (m Member) GetTiers() []models.Tier {
	return m.store.GetTiers()
}

func (m Member) FindNonMembersOnSlack() []string {
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
