package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"memberserver/api/models"
	"memberserver/database"
	"memberserver/payments"
	"memberserver/slack"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (a API) getTiers(w http.ResponseWriter, req *http.Request) {
	tiers := a.db.GetMemberTiers()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(tiers)
	w.Write(j)
}

func (a API) getMembers(w http.ResponseWriter, req *http.Request) {
	members := a.db.GetMembers()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(members)
	w.Write(j)
}

func (a API) getMemberByEmail(w http.ResponseWriter, req *http.Request) {
	routeVars := mux.Vars(req)

	memberEmail := routeVars["email"]

	member, err := a.db.GetMemberByEmail((memberEmail))

	if err != nil {
		log.Errorf("error getting member by email: %s", err)
		http.Error(w, errors.New("error getting member by email").Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(member)
	w.Write(j)
}

func (a API) assignRFID(w http.ResponseWriter, req *http.Request) {
	var assignRFIDRequest database.AssignRFIDRequest

	err := json.NewDecoder(req.Body).Decode(&assignRFIDRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r, err := a.db.SetRFIDTag(assignRFIDRequest.Email, assignRFIDRequest.RFID)
	if err != nil {
		log.Errorf("error trying to assign rfid to member: %s", err.Error())
		http.Error(w, errors.New("unable to assign rfid").Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(r)
	w.Write(j)
}

func (a API) refreshPayments(w http.ResponseWriter, req *http.Request) {
	payments.GetPayments()

	a.db.EvaluateMembers()

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(models.EndpointSuccess{
		Ack: true,
	})
	w.Write(j)
}

func (a API) getNonMembersOnSlack(w http.ResponseWriter, req *http.Request) {
	nonMembers := slack.FindNonMembers()
	buf := bytes.NewBufferString(strings.Join(nonMembers[:], "\n"))
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=nonmembersOnSlack.csv")
	w.Write(buf.Bytes())
}
