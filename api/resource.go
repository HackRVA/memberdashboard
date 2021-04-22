package api

import (
	"encoding/json"
	"memberserver/api/models"
	"memberserver/database"
	"net/http"

	"memberserver/resourcemanager"

	log "github.com/sirupsen/logrus"
)

// Resource http handlers for resources
func (rs resourceAPI) Resource(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		rs.get(w, req)
	}

	if req.Method == http.MethodPut {
		rs.update(w, req)
	}

	if req.Method == http.MethodDelete {
		rs.delete(w, req)
	}
}

func (rs resourceAPI) get(w http.ResponseWriter, req *http.Request) {
	resourceList := rs.db.GetResources()

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(resourceList)
	w.Write(j)
}

func (rs resourceAPI) update(w http.ResponseWriter, req *http.Request) {
	var updateResourceReq database.Resource

	err := json.NewDecoder(req.Body).Decode(&updateResourceReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r, err := rs.db.UpdateResource(updateResourceReq.ID, updateResourceReq.Name, updateResourceReq.Address, updateResourceReq.IsDefault)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(r)
	w.Write(j)
}

func (rs resourceAPI) delete(w http.ResponseWriter, req *http.Request) {
	var deleteResourceReq database.ResourceDeleteRequest

	err := json.NewDecoder(req.Body).Decode(&deleteResourceReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("attempting to delete %s", deleteResourceReq.ID)

	err = rs.db.DeleteResource(deleteResourceReq.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(models.EndpointSuccess{
		Ack: true,
	})
	w.Write(j)
}

func (rs resourceAPI) addMultipleMembersToResource(w http.ResponseWriter, req *http.Request) {
	var membersResource models.MembersResourceRelation

	err := json.NewDecoder(req.Body).Decode(&membersResource)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resource, err := rs.db.AddMultipleMembersToResource(membersResource.Emails, membersResource.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(resource)
	w.Write(j)
}

func (rs resourceAPI) removeMember(w http.ResponseWriter, req *http.Request) {
	var update models.MemberResourceRelation

	err := json.NewDecoder(req.Body).Decode(&update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rs.db.RemoveUserFromResource(update.Email, update.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(models.EndpointSuccess{
		Ack: true,
	})
	w.Write(j)

	resource, err := rs.db.GetResourceByID(update.ID)
	if err != nil {
		log.Errorf("error getting resource to update when removing a member: %s", err)
	}

	resourcemanager.UpdateResourceACL(resource)
	resourcemanager.UpdateResources()
}

func (rs resourceAPI) register(w http.ResponseWriter, req *http.Request) {
	var register database.RegisterResourceRequest

	err := json.NewDecoder(req.Body).Decode(&register)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r, err := rs.db.RegisterResource(register.Name, register.Address, register.IsDefault)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(r)
	w.Write(j)
}

func (rs resourceAPI) status(w http.ResponseWriter, req *http.Request) {
	resources := rs.db.GetResources()
	// statusMap := make(map[string]uint8)

	for _, r := range resources {
		if r == (database.Resource{}) {
			continue
		}
		resourcemanager.CheckStatus(r)
		// if err != nil {
		// 	log.Errorf("error getting resource status: %s", err.Error())
		// 	statusMap[r.Name] = resourcemanager.StatusOffline
		// 	continue
		// }
		// if status == resourcemanager.StatusOutOfDate {
		// 	statusMap[r.Name] = resourcemanager.StatusOutOfDate
		// 	continue
		// }
		// statusMap[r.Name] = resourcemanager.StatusGood
	}

	w.Header().Set("Content-Type", "application/json")
	// j, _ := json.Marshal(statusMap)
	// w.Write(j)
}

func (rs resourceAPI) updateResourceACL(w http.ResponseWriter, req *http.Request) {
	resourcemanager.UpdateResources()

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(models.EndpointSuccess{
		Ack: true,
	})
	w.Write(j)
}

func (rs resourceAPI) deleteResourceACL(w http.ResponseWriter, req *http.Request) {
	resourcemanager.DeleteResourceACL()

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(models.EndpointSuccess{
		Ack: true,
	})
	w.Write(j)
}
