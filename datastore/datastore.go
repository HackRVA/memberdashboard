package datastore

import (
	"context"
	"time"

	"github.com/HackRVA/memberserver/models"
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
		LogAccessEvent(ctx context.Context, event models.LogMessage) error
	}

	MemberStore interface {
		GetTiers(ctx context.Context) []models.Tier
		GetMembers(ctx context.Context) []models.Member
		GetMemberCount(ctx context.Context, isActive bool) (int, error)
		GetMembersPaginated(ctx context.Context, limit int, page int, active bool) ([]models.Member, error)
		GetMemberByEmail(ctx context.Context, email string) (models.Member, error)
		AssignRFID(ctx context.Context, email string, rfid string) (models.Member, error)
		AddNewMember(ctx context.Context, newMember models.Member) (models.Member, error)
		GetMembersWithCredit(ctx context.Context) []models.Member
		ProcessMember(ctx context.Context, newMember models.Member) error
		GetMemberByRFID(ctx context.Context, rfid string) (models.Member, error)
		UpdateMember(ctx context.Context, update models.Member) error
		UpdateMemberByID(ctx context.Context, memberID string, update models.Member) error
		UpdateMemberBySubscriptionID(ctx context.Context, subscriptionID string, update models.Member) error
		SetMemberLevel(ctx context.Context, memberId string, level models.MemberLevel) error
		ApplyMemberCredits(ctx context.Context)
		UpdateMemberTiers(ctx context.Context)
		GetActiveMembersWithoutSubscription(ctx context.Context) []models.Member
	}

	ResourceStore interface {
		GetResources(ctx context.Context) []models.Resource
		GetResourceByID(ctx context.Context, ID string) (models.Resource, error)
		GetResourceByName(ctx context.Context, resourceName string) (models.Resource, error)
		RegisterResource(ctx context.Context, name string, address string, isDefault bool) (models.Resource, error)
		UpdateResource(ctx context.Context, res models.Resource) (*models.Resource, error)
		DeleteResource(ctx context.Context, id string) error
		AddMultipleMembersToResource(ctx context.Context, emails []string, resourceID string) ([]models.MemberResourceRelation, error)
		AddUserToDefaultResources(ctx context.Context, email string) ([]models.MemberResourceRelation, error)
		RemoveUserFromResource(ctx context.Context, email string, resourceID string) error
		GetResourceACL(ctx context.Context, r models.Resource) ([]string, error)
		GetResourceACLWithMemberInfo(ctx context.Context, r models.Resource) ([]models.Member, error)
		GetMembersAccess(ctx context.Context, m models.Member) ([]models.MemberAccess, error)
		GetInactiveMembersByResource(ctx context.Context) ([]models.MemberAccess, error)
		GetActiveMembersByResource(ctx context.Context) ([]models.MemberAccess, error)
	}

	CommunicationStore interface {
		GetCommunications(ctx context.Context) []models.Communication
		GetCommunication(ctx context.Context, name string) (models.Communication, error)
		GetMostRecentCommunicationToMember(ctx context.Context, memberId string, commId int) (time.Time, error)
		LogCommunication(ctx context.Context, communicationId int, memberId string) error
	}

	UserStore interface {
		GetUser(ctx context.Context, email string) (models.UserResponse, error)
		UserSignin(ctx context.Context, email string, password string) error
		RegisterUser(ctx context.Context, creds models.Credentials) error
	}

	ReportStore interface {
		UpdateMemberCounts(ctx context.Context)
		GetMemberCounts(ctx context.Context) ([]models.MemberCount, error)
		GetMemberCountByMonth(ctx context.Context, month time.Time) (models.MemberCount, error)
		GetAccessStats(ctx context.Context, date time.Time, resourceName string) ([]models.AccessStats, error)
		GetMemberChurn(ctx context.Context) (int, error)
	}
)
