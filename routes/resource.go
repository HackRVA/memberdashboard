package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

// ResourceRequest to update or delete a resource
type ResourceRequest struct {
	ID      uint8  `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type registerResourceRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Resource http handlers for resources
func (a *API) Resource(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		a.getResources(w, req)
	}

	if req.Method == http.MethodPost {
		a.updateResource(w, req)
	}

	if req.Method == http.MethodDelete {
		a.deleteResource(w, req)
	}
}

func (a *API) getResources(w http.ResponseWriter, req *http.Request) {
	resourceList := a.db.GetResources()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(resourceList)
	w.Write(j)
}

func (a *API) updateResource(w http.ResponseWriter, req *http.Request) {
	var updateResourceReq ResourceRequest

	err := json.NewDecoder(req.Body).Decode(&updateResourceReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r, err := a.db.UpdateResource(updateResourceReq.ID, updateResourceReq.Name, updateResourceReq.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(r)
	w.Write(j)
}

func (a *API) deleteResource(w http.ResponseWriter, req *http.Request) {
	var deleteResourceReq ResourceRequest

	err := json.NewDecoder(req.Body).Decode(&deleteResourceReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("attempting to delete %s", deleteResourceReq.Name)

	err = a.db.DeleteResource(deleteResourceReq.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
