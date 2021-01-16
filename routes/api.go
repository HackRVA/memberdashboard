package routes

import (
	"encoding/json"

	"github.com/dfirebaugh/memberserver/config"
	"github.com/dfirebaugh/memberserver/database"

	"net/http"
)

type resourceAPI struct {
	db     *database.Database
	config config.Config
}

// API endpoints
type API struct {
	db       *database.Database
	config   config.Config
	resource resourceAPI
}

func (a *API) info(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(struct{ Message string }{
		Message: "hello, world!",
	})
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
