package api

import (
	"encoding/json"
	"memberserver/database"
	"net/http"
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

func (a *API) getTiers(w http.ResponseWriter, req *http.Request) {
	tiers := a.db.GetMemberTiers()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(tiers)
	w.Write(j)
}

func (a *API) getMembers(w http.ResponseWriter, req *http.Request) {
	members := a.db.GetMembers()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(members)
	w.Write(j)
}
