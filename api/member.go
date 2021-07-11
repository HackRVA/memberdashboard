package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"memberserver/api/models"
	"memberserver/resourcemanager"
	"memberserver/slack"
	"net/http"
	"strings"
)

type MemberStore interface {
	GetTiers() []models.Tier // update where this is
	GetMembers() []models.Member
	GetMemberByEmail(email string) (models.Member, error)
	AssignRFID(email string, rfid string) (models.Member, error)
	AddNewMember(newMember models.Member) (models.Member, error)
}

type MemberServer struct {
	store MemberStore
}

func (m *MemberServer) GetMembersHandler(w http.ResponseWriter, r *http.Request) {
	members := m.store.GetMembers()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(members)
	w.Write(j)
}

func (m *MemberServer) GetByEmailHandler(w http.ResponseWriter, r *http.Request) {
	memberEmail := strings.TrimPrefix(r.URL.Path, "/api/member/email/")

	if len(memberEmail) == 0 {
		http.Error(w, errors.New("error getting member by email").Error(), http.StatusNotFound)
		return
	}

	member, err := m.store.GetMemberByEmail(memberEmail)

	if err != nil {
		http.Error(w, "error getting member by email", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(member)
	w.Write(j)
}

func (m *MemberServer) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	_, user, _ := strategy.AuthenticateRequest(r)

	member, err := m.store.GetMemberByEmail(user.GetUserName())

	if err != nil {
		http.Error(w, "error getting member by email", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(member)
	w.Write(j)
}

func (m *MemberServer) AssignRFIDHandler(w http.ResponseWriter, r *http.Request) {
	var assignRFIDRequest models.AssignRFIDRequest

	err := json.NewDecoder(r.Body).Decode(&assignRFIDRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	_, user, _ := strategy.AuthenticateRequest(r)

	m.assignRFID(w, user.GetUserName(), assignRFIDRequest.RFID)
}

func (m *MemberServer) GetTiersHandler(w http.ResponseWriter, r *http.Request) {
	tiers := m.store.GetTiers()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(tiers)
	w.Write(j)
}

func (m *MemberServer) assignRFID(w http.ResponseWriter, email, rfid string) {
	if len(rfid) == 0 {
		http.Error(w, errors.New("not a valid rfid").Error(), http.StatusBadRequest)
		return
	}
	r, err := m.store.AssignRFID(email, rfid)
	if err != nil {
		http.Error(w, errors.New("unable to assign rfid").Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(r)
	w.Write(j)

	go resourcemanager.PushOne(models.Member{Email: email})
}

func (m *MemberServer) GetNonMembersOnSlackHandler(w http.ResponseWriter, r *http.Request) {
	nonMembers := slack.FindNonMembers()
	buf := bytes.NewBufferString(strings.Join(nonMembers[:], "\n"))
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=nonmembersOnSlack.csv")
	w.Write(buf.Bytes())
}

func (m *MemberServer) AddNewMemberHandler(w http.ResponseWriter, r *http.Request) {
	var newMember models.Member

	err := json.NewDecoder(r.Body).Decode(&newMember)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addedMember, err := m.store.AddNewMember(newMember)
	if err != nil {
		http.Error(w, "error getting member by email", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(addedMember)
	w.Write(j)

	m.store.AssignRFID(addedMember.Email, addedMember.RFID)

	go resourcemanager.PushOne(addedMember)
}
