package v1

import (
	"net/http"

	"github.com/HackRVA/memberserver/pkg/paypal/listener"
)

type V1Router interface {
	SetupRoutes()
	MemberHTTPHandler
	ReportsHTTPHandler
	ResourceHTTPHandler
	UserHTTPHandler
	AuthHTTPHandler
	VersionHTTPHandler
	PaymentHTTPHandler
}

type Logger interface {
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Tracef(format string, args ...interface{})
	Print(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Trace(args ...interface{})
}

type MemberHTTPHandler interface {
	MemberEmail(w http.ResponseWriter, r *http.Request)
	GetMembers(w http.ResponseWriter, r *http.Request)
	UpdateMemberByEmail(w http.ResponseWriter, r *http.Request)
	GetByEmail(w http.ResponseWriter, r *http.Request)
	GetCurrentUser(w http.ResponseWriter, r *http.Request)
	AssignRFID(w http.ResponseWriter, r *http.Request)
	AssignRFIDSelf(w http.ResponseWriter, r *http.Request)
	GetTiers(w http.ResponseWriter, r *http.Request)
	GetNonMembersOnSlack(w http.ResponseWriter, r *http.Request)
	AddNewMember(w http.ResponseWriter, r *http.Request)
	CheckStatus(w http.ResponseWriter, r *http.Request)
	SetCredited(w http.ResponseWriter, r *http.Request)
}

type ReportsHTTPHandler interface {
	GetMemberCountsCharts(http.ResponseWriter, *http.Request)
	GetAccessStatsChart(http.ResponseWriter, *http.Request)
	GetMemberChurn(http.ResponseWriter, *http.Request)
}

type ResourceHTTPHandler interface {
	Resource(w http.ResponseWriter, req *http.Request)
	AddMultipleMembersToResource(w http.ResponseWriter, req *http.Request)
	RemoveMember(w http.ResponseWriter, req *http.Request)
	Register(w http.ResponseWriter, req *http.Request)
	Status(w http.ResponseWriter, req *http.Request)
	UpdateResourceACL(w http.ResponseWriter, req *http.Request)
	Open(w http.ResponseWriter, req *http.Request)
	DeleteResourceACL(w http.ResponseWriter, req *http.Request)
}

type UserHTTPHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
}

type AuthHTTPHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type VersionHTTPHandler interface {
	Version(http.ResponseWriter, *http.Request)
}

type PaymentHTTPHandler interface {
	PaypalSubscriptionWebHook(err error, n *listener.Subscription)
}
