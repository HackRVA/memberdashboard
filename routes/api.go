package routes

import (
	"encoding/json"

	"github.com/dfirebaugh/memberserver/database"

	"net/http"
)

// API endpoints
type API struct {
	db *database.Database
}

func (a *API) info(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(struct{ Message string }{
		Message: "hello, world!",
	})
	w.Write(j)
}

func (a *API) getResources(w http.ResponseWriter, req *http.Request) {
	resourceList := a.db.GetResources()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(resourceList)
	w.Write(j)
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

type registerResourceRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (a *API) registerResource(w http.ResponseWriter, req *http.Request) {
	var newResourceReq registerResourceRequest

	err := json.NewDecoder(req.Body).Decode(&newResourceReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r, err := a.db.RegisterResource(newResourceReq.Name, newResourceReq.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(r)
	w.Write(j)
}
