package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/HackRVA/memberserver/models"
	"github.com/HackRVA/memberserver/services"
	"github.com/sirupsen/logrus"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
)

type MemberServer struct {
	ResourceManager services.Resource
	MemberService   services.Member
	AuthStrategy    union.Union
}

func (m *MemberServer) MemberEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		m.GetByEmail(w, r)
	}

	if r.Method == http.MethodPut {
		m.UpdateMemberByEmail(w, r)
	}
}

func (m *MemberServer) GetMembers(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search != "" {
		results := []models.Member{}
		for _, member := range m.MemberService.Get() {
			if fuzzy.Match(strings.ToLower(search), strings.ToLower(member.Name)) ||
				fuzzy.Match(strings.ToLower(search), strings.ToLower(member.Email)) ||
				fuzzy.Match(strings.ToLower(search), strings.ToLower(member.RFID)) ||
				fuzzy.Match(strings.ToLower(search), strings.ToLower(member.SubscriptionID)) {
				results = append(results, member)
			}
		}
		limit := 10
		if len(results) > limit {
			results = results[:limit]
		}
		ok(w, models.MembersPaginatedResponse{
			Members: results,
			Count:   uint(len(results)),
		})
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

	getActive := active == "true"

	memberCount, err := m.MemberService.GetMemberCount(getActive)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	ok(w, models.MembersPaginatedResponse{
		Members: m.MemberService.GetMembersPaginated(count, page, active == "true"),
		Count:   uint(memberCount),
	})
}

func (m *MemberServer) UpdateMemberByEmail(w http.ResponseWriter, r *http.Request) {
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
		notFound(w, fmt.Sprintf("error updating member: %s", err))
		return
	}

	ok(w, models.EndpointSuccess{
		Ack: true,
	})
}

func (m *MemberServer) GetByEmail(w http.ResponseWriter, r *http.Request) {
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

func (m *MemberServer) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	_, user, _ := m.AuthStrategy.AuthenticateRequest(r)

	member, err := m.MemberService.GetByEmail(user.GetUserName())
	if err != nil {
		notFound(w, "error getting member by email")
		return
	}

	ok(w, member)
}

func (m *MemberServer) AssignRFID(w http.ResponseWriter, r *http.Request) {
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

func (m *MemberServer) AssignRFIDSelf(w http.ResponseWriter, r *http.Request) {
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

func (m *MemberServer) GetTiers(w http.ResponseWriter, r *http.Request) {
	ok(w, m.MemberService.GetTiers())
}

func (m *MemberServer) GetNonMembersOnSlack(w http.ResponseWriter, r *http.Request) {
	nonMembers := m.MemberService.FindNonMembersOnSlack()
	buf := bytes.NewBufferString(strings.Join(nonMembers[:], "\n"))
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=nonmembersOnSlack.csv")
	if _, err := w.Write(buf.Bytes()); err != nil {
		logrus.Error(err)
	}
}

func (m *MemberServer) AddNewMember(w http.ResponseWriter, r *http.Request) {
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

func (m *MemberServer) CheckStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		badRequest(w, "not a valid subscriptionID")
		return
	}

	member, err := m.MemberService.CheckStatus(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting member by status: %s", err.Error()), http.StatusNotFound)
		return
	}

	ok(w, member)
}

func (m *MemberServer) SetCredited(w http.ResponseWriter, r *http.Request) {
	var creditRequest models.MemberShipCreditRequest
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		badRequest(w, "not a valid subscriptionID")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&creditRequest)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	level := models.Standard

	if creditRequest.IsCredited {
		level = models.Credited
	}

	err = m.MemberService.SetLevel(id, level)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting member by status: %s", err.Error()), http.StatusNotFound)
		return
	}

	ok(w, models.InfoResponse{
		Message: fmt.Sprintf("member %s set to level %s", id, models.MemberLevelToStr[level]),
	})
}
