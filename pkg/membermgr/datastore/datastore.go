package datastore

import (
	"time"

	"github.com/HackRVA/memberserver/pkg/membermgr/models"
)

type (
	DataStore interface {
		AccessEvent
		MemberStore
		ResourceStore
		CommunicationStore
		UserStore
		ReportStore
	}

	AccessEvent interface {
		LogAccessEvent(event models.LogMessage) error
	}

	MemberStore interface {
		GetTiers() []models.Tier // update where this is
		GetMembers() []models.Member
		GetMembersWithLimit(limit int, offset int, active bool) []models.Member
		GetMemberByEmail(email string) (models.Member, error)
		AssignRFID(email string, rfid string) (models.Member, error)
		AddNewMember(newMember models.Member) (models.Member, error)
		AddMembers(members []models.Member) error
		GetMembersWithCredit() []models.Member
		ProcessMember(newMember models.Member) error
		GetMemberByRFID(rfid string) (models.Member, error)
		UpdateMember(update models.Member) error
		UpdateMemberBySubscriptionID(subscriptionID string, update models.Member) error
		SetMemberLevel(memberId string, level models.MemberLevel) error
		ApplyMemberCredits()
		UpdateMemberTiers()
	}

	ResourceStore interface {
		GetResources() []models.Resource
		GetResourceByID(ID string) (models.Resource, error)
		GetResourceByName(resourceName string) (models.Resource, error)
		RegisterResource(name string, address string, isDefault bool) (models.Resource, error)
		UpdateResource(res models.Resource) (*models.Resource, error)
		DeleteResource(id string) error
		AddMultipleMembersToResource(emails []string, resourceID string) ([]models.MemberResourceRelation, error)
		AddUserToDefaultResources(email string) ([]models.MemberResourceRelation, error)
		GetMemberResourceRelation(m models.Member, r models.Resource) (models.MemberResourceRelation, error)
		RemoveUserFromResource(email string, resourceID string) error
		GetResourceACL(r models.Resource) ([]string, error)
		GetResourceACLWithMemberInfo(r models.Resource) ([]models.Member, error)
		GetMembersAccess(m models.Member) ([]models.MemberAccess, error)
		GetInactiveMembersByResource() ([]models.MemberAccess, error)
		GetActiveMembersByResource() ([]models.MemberAccess, error)
	}

	CommunicationStore interface {
		GetCommunications() []models.Communication
		GetCommunication(name string) (models.Communication, error)
		GetMostRecentCommunicationToMember(memberId string, commId int) (time.Time, error)
		LogCommunication(communicationId int, memberId string) error
	}

	UserStore interface {
		GetUser(email string) (models.UserResponse, error)
		UserSignin(email string, password string) error
		RegisterUser(creds models.Credentials) error
	}

	ReportStore interface {
		UpdateMemberCounts()
		GetMemberCounts() ([]models.MemberCount, error)
		GetMemberCountByMonth(month time.Time) (models.MemberCount, error)
		GetAccessStats(date time.Time, resourceName string) ([]models.AccessStats, error)
		GetMemberChurn() (int, error)
	}
)
