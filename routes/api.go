package routes

import (
	"encoding/json"

	"github.com/dfirebaugh/memberserver/database"

	"net/http"
)

// API endpoints
type API struct {
	DB *database.Database
}

func (a *API) info(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(struct{ Message string }{
		Message: "hello, world!",
	})
	w.Write(j)
}

func (a *API) getStatuses(w http.ResponseWriter, req *http.Request) {
	statusList := a.DB.GetStatuses()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(statusList)
	w.Write(j)
}

func (a *API) getResources(w http.ResponseWriter, req *http.Request) {
	resourceList := a.DB.GetResources()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(resourceList)
	w.Write(j)
}

func (a *API) getTiers(w http.ResponseWriter, req *http.Request) {
	tiers := a.DB.GetMemberTiers()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(tiers)
	w.Write(j)
}

func (a *API) getMembers(w http.ResponseWriter, req *http.Request) {
	members := a.DB.GetMembers()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(members)
	w.Write(j)
}
