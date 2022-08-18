package controllers

import (
	"bytes"
	"encoding/json"
	"memberserver/internal/models"
	"memberserver/internal/services/member"
	"memberserver/internal/services/resourcemanager"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
)

type MemberServer struct {
	ResourceManager resourcemanager.ResourceManager
	MemberService   member.MemberService
	AuthStrategy    union.Union
}

func (m *MemberServer) MemberEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		m.GetByEmailHandler(w, r)
	}

	if r.Method == http.MethodPut {
		m.UpdateMemberByEmailHandler(w, r)
	}
}

func (m *MemberServer) GetMembersHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search != "" {
		results := []models.Member{}
		for _, member := range m.MemberService.Get() {
			if fuzzy.Match(search, member.Name) || fuzzy.Match(search, member.Email) || fuzzy.Match(search, member.RFID) || fuzzy.Match(search, member.SubscriptionID) {
				results = append(results, member)
			}
		}
		ok(w, results)
		return
	}
	active := r.URL.Query().Get("active")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		ok(w, m.MemberService.Get())
		return
	}
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		ok(w, m.MemberService.Get())
		return
	}

	if active == "true" {
		println("get active")
	}

	ok(w, m.MemberService.GetMembersWithLimit(count, page, active == "true"))
}

func (m *MemberServer) UpdateMemberByEmailHandler(w http.ResponseWriter, r *http.Request) {
	memberEmail := strings.TrimPrefix(r.URL.Path, "/api/member/email/")

	var request models.UpdateMemberRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if len(memberEmail) == 0 || !govalidator.IsEmail(memberEmail) {
		preconditionFailed(w, "invalid email")
		return
	}

	if err != nil {
		badRequest(w, err.Error())
		return
	}

	if len(request.FullName) == 0 {
		preconditionFailed(w, "fullName is required")
		return
	}

	err = m.MemberService.Update(models.Member{
		Email:          memberEmail,
		Name:           request.FullName,
		SubscriptionID: request.SubscriptionID,
	})

	if err != nil {
		notFound(w, "error getting member by email")
		return
	}

	ok(w, models.EndpointSuccess{
		Ack: true,
	})

}

func (m *MemberServer) GetByEmailHandler(w http.ResponseWriter, r *http.Request) {
	memberEmail := strings.TrimPrefix(r.URL.Path, "/api/member/email/")

	if len(memberEmail) == 0 || !govalidator.IsEmail(memberEmail) {
		preconditionFailed(w, "invalid email")
		return
	}

	member, err := m.MemberService.GetByEmail(memberEmail)

	if err != nil {
		notFound(w, "error getting member by email")
		return
	}

	ok(w, member)
}

func (m *MemberServer) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	_, user, _ := m.AuthStrategy.AuthenticateRequest(r)

	member, err := m.MemberService.GetByEmail(user.GetUserName())

	if err != nil {
		notFound(w, "error getting member by email")
		return
	}

	ok(w, member)
}

func (m *MemberServer) AssignRFIDHandler(w http.ResponseWriter, r *http.Request) {
	var assignRFIDRequest models.AssignRFIDRequest

	err := json.NewDecoder(r.Body).Decode(&assignRFIDRequest)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	if len(assignRFIDRequest.RFID) == 0 {
		preconditionFailed(w, "not a valid rfid")
		return
	}

	member, err := m.MemberService.AssignRFID(assignRFIDRequest.Email, assignRFIDRequest.RFID)
	if err != nil {
		notFound(w, "unable to assign rfid")
		return
	}

	ok(w, member)
}

func (m *MemberServer) AssignRFIDSelfHandler(w http.ResponseWriter, r *http.Request) {
	var assignRFIDRequest models.AssignRFIDRequest

	err := json.NewDecoder(r.Body).Decode(&assignRFIDRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, user, _ := m.AuthStrategy.AuthenticateRequest(r)

	member, err := m.MemberService.AssignRFID(user.GetUserName(), assignRFIDRequest.RFID)
	if err != nil {
		notFound(w, "unable to assign rfid")
		return
	}

	ok(w, member)
}

func (m *MemberServer) GetTiersHandler(w http.ResponseWriter, r *http.Request) {
	ok(w, m.MemberService.GetTiers())
}

func (m *MemberServer) GetNonMembersOnSlackHandler(w http.ResponseWriter, r *http.Request) {
	nonMembers := m.MemberService.FindNonMembersOnSlack()
	buf := bytes.NewBufferString(strings.Join(nonMembers[:], "\n"))
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=nonmembersOnSlack.csv")
	w.Write(buf.Bytes())
}

func (m *MemberServer) AddNewMemberHandler(w http.ResponseWriter, r *http.Request) {
	var newMember models.Member

	err := json.NewDecoder(r.Body).Decode(&newMember)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	addedMember, err := m.MemberService.Add(newMember)
	if err != nil {
		http.Error(w, "error getting member by email", http.StatusNotFound)
		return
	}

	ok(w, addedMember)

}
