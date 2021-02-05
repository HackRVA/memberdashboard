package api

import (
	"encoding/json"
	"errors"
	"memberserver/database"
	"memberserver/payments"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// swagger:response getMemberResponse
type memberResponseBody struct {
	// in: body
	Body []database.Member
}

// swagger:response getTierResponse
type getTierResponse struct {
	// in: body
	Body []database.Tier
}

// swagger:response setRFIDResponse
type setRFIDResponse struct {
	// in: body
	Body database.AssignRFIDRequest
}

// swagger:parameters setRFIDRequest
type setRFIDRequest struct {
	// in: body
	Body database.AssignRFIDRequest
}

// swagger:response getPaymentRefreshResponse
type getPaymentRefreshResponse struct {
	Body endpointSuccess
}

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
	members := a.db.GetMembers()

	for _, m := range members {
		err := a.db.EvaluateMemberStatus(m.ID)
		if err != nil {
			log.Errorf("error evaluating member's status: %s", err.Error())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(endpointSuccess{
		Ack: true,
	})
	w.Write(j)
}
