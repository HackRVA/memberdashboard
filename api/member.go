package api

import (
	"bytes"
	"encoding/json"
	"memberserver/api/models"
	"memberserver/datastore"
	"memberserver/resourcemanager"
	"memberserver/slack"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	log "github.com/sirupsen/logrus"
)

type MemberServer struct {
	store           datastore.DataStore
	ResourceManager *resourcemanager.ResourceManager
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
	members := m.store.GetMembers()

	ok(w, members)
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

	_, err = m.store.GetMemberByEmail(memberEmail)

	if err != nil {
		notFound(w, "error getting member by email")
		return
	}

	err = m.store.UpdateMemberByEmail(request.FullName, memberEmail)

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

	member, err := m.store.GetMemberByEmail(memberEmail)

	if err != nil {
		notFound(w, "error getting member by email")
		return
	}

	ok(w, member)
}

func (m *MemberServer) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	_, user, _ := m.AuthStrategy.AuthenticateRequest(r)

	member, err := m.store.GetMemberByEmail(user.GetUserName())

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

	m.assignRFID(w, assignRFIDRequest.Email, assignRFIDRequest.RFID)
}

func (m *MemberServer) AssignRFIDSelfHandler(w http.ResponseWriter, r *http.Request) {
	var assignRFIDRequest models.AssignRFIDRequest

	err := json.NewDecoder(r.Body).Decode(&assignRFIDRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, user, _ := m.AuthStrategy.AuthenticateRequest(r)

	m.assignRFID(w, user.GetUserName(), assignRFIDRequest.RFID)
}

func (m *MemberServer) GetTiersHandler(w http.ResponseWriter, r *http.Request) {
	tiers := m.store.GetTiers()

	ok(w, tiers)
}

func (m *MemberServer) assignRFID(w http.ResponseWriter, email, rfid string) {
	if len(rfid) == 0 {
		preconditionFailed(w, "not a valid rfid")
		return
	}

	m.removeMembersRFID(email)

	r, err := m.store.AssignRFID(email, rfid)
	if err != nil {
		notFound(w, "unable to assign rfid")
		return
	}

	ok(w, r)

	go m.ResourceManager.PushOne(models.Member{Email: email})
}

func (m *MemberServer) removeMembersRFID(email string) {
	member, err := m.store.GetMemberByEmail(email)
	if err != nil {
		log.Error(err)
		return
	}

	if member.RFID == "notset" || len(member.RFID) > 0 {
		return
	}

	for _, r := range member.Resources {
		resource, err := m.store.GetResourceByID(r.ResourceID)
		if err != nil {
			log.Error(err)
			continue
		}

		m.ResourceManager.RemoveMember(models.MemberAccess{
			Email:           member.Email,
			ResourceAddress: resource.Address,
			ResourceName:    resource.Name,
			Name:            member.Name,
			RFID:            member.RFID,
		})

	}
}

func (m *MemberServer) GetNonMembersOnSlackHandler(w http.ResponseWriter, r *http.Request) {
	nonMembers := slack.FindNonMembers(m.store)
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

	addedMember, err := m.store.AddNewMember(newMember)
	if err != nil {
		http.Error(w, "error getting member by email", http.StatusNotFound)
		return
	}

	ok(w, addedMember)

	m.store.AssignRFID(addedMember.Email, addedMember.RFID)

	go m.ResourceManager.PushOne(addedMember)
}
