package routes

import (
	"encoding/json"

	"github.com/dfirebaugh/memberserver/database"

	"net/http"
)

func info(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(struct{ Message string }{
		Message: "hello, world!",
	})
	w.Write(j)
}

func getStatuses(w http.ResponseWriter, req *http.Request) {
	statusList := database.DB.GetStatuses()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(statusList)
	w.Write(j)
}

func getResources(w http.ResponseWriter, req *http.Request) {
	resourceList := database.DB.GetResources()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(resourceList)
	w.Write(j)
}

func getTiers(w http.ResponseWriter, req *http.Request) {
	tiers := database.DB.GetMemberTiers()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(tiers)
	w.Write(j)
}

func getMembers(w http.ResponseWriter, req *http.Request) {
	members := database.DB.GetMembers()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(members)
	w.Write(j)
}
