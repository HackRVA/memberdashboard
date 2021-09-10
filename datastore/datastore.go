package datastore

import (
	"memberserver/api/models"
	"time"
)

type DataStore interface {
	Event
	MemberStore
	ResourceStore
	PaymentStore
	CommunicationStore
	UserStore
}

type Event interface {
	AddLogMsg([]byte) error
}

type MemberStore interface {
	GetTiers() []models.Tier // update where this is
	GetMembers() []models.Member
	GetMemberByEmail(email string) (models.Member, error)
	AssignRFID(email string, rfid string) (models.Member, error)
	AddNewMember(newMember models.Member) (models.Member, error)
	AddMembers(members []models.Member) error
	GetMembersWithCredit() []models.Member
	ProcessMember(newMember models.Member) error
	GetMemberByRFID(rfid string) (models.Member, error)
}

type ResourceStore interface {
	GetResources() []models.Resource
	GetResourceByID(ID string) (models.Resource, error)
	GetResourceByName(resourceName string) (models.Resource, error)
	RegisterResource(name string, address string, isDefault bool) (models.Resource, error)
	UpdateResource(id string, name string, address string, isDefault bool) (*models.Resource, error)
	DeleteResource(id string) error
	AddMultipleMembersToResource(emails []string, resourceID string) ([]models.MemberResourceRelation, error)
	AddUserToDefaultResources(email string) ([]models.MemberResourceRelation, error)
	GetMemberResourceRelation(m models.Member, r models.Resource) (models.MemberResourceRelation, error)
	RemoveUserFromResource(email string, resourceID string) error
	GetResourceACL(r models.Resource) ([]string, error)
	GetResourceACLWithMemberInfo(r models.Resource) ([]models.Member, error)
	GetMembersAccess(m models.Member) ([]models.MemberAccess, error)
}

type PaymentStore interface {
	GetPayments() ([]models.Payment, error)
	AddPayment(payment models.Payment) error
	AddPayments(payments []models.Payment) error
	SetMemberLevel(memberId string, level models.MemberLevel) error
	ApplyMemberCredits()
	UpdateMemberTiers()
	GetPastDueAccounts() []models.PastDueAccount
}

type CommunicationStore interface {
	GetCommunications() []models.Communication
	GetCommunication(name string) (models.Communication, error)
	GetMostRecentCommunicationToMember(memberId string, commId int) (time.Time, error)
	LogCommunication(communicationId int, memberId string) error
}

type UserStore interface {
	GetUser(email string) (models.UserResponse, error)
	UserSignin(email string, password string) error
	RegisterUser(creds models.Credentials) error
}
