package api

import (
	"encoding/json"
	"memberserver/database"
	"net/http"

	"memberserver/resourcemanager"

	log "github.com/sirupsen/logrus"
)

// resource to update or delete a resource
type memberResourceRelation struct {
	// ID of the Resource
	// required: true
	// example: string
	ID string `json:"resourceID"`
	// Email - this will be the member's email address
	// Name of the Resource
	// required: true
	// example: email
	Email string `json:"email"`
}

type membersResourceRelation struct {
	// ID of the Resource
	// required: true
	// example: string
	ID string `json:"resourceID"`
	// Emails - list of member's email address
	// required: true
	// example: []
	Emails []string `json:"emails"`
}

// swagger:parameters updateResourceRequest
type updateResourceRequest struct {
	// in: body
	Body database.ResourceRequest
}

// swagger:parameters registerResourceRequest
type registerResourceRequest struct {
	// in: body
	Body database.RegisterResourceRequest
}

// swagger:parameters deleteResourceRequest
type deleteResourceRequest struct {
	// in: body
	Body database.ResourceDeleteRequest
}

// swagger:parameters resourceAddMemberRequest
type resourceAddMemberRequest struct {
	// in: body
	Body memberResourceRelation
}

// swagger:parameters resourceBulkMemberRequest
type resourceBulkMemberRequest struct {
	// in: body
	Body membersResourceRelation
}

// swagger:response getResourceResponse
type getResourceResponse struct {
	// in: body
	Body database.Resource
}

// swagger:response postResourceResponse
type postResourceResponse struct {
	// in: body
	Body database.Resource
}

// swagger:response addMemberToResourceResponse
type addMemberToResourceResponse struct {
	// in: body
	Body database.MemberResourceRelation
}

// swagger:response addMulitpleMembersToResourceResponse
type addMulitpleMembersToResourceResponse struct {
	// in: body
	Body []database.MemberResourceRelation
}

// swagger:response getResourceStatusResponse
type getResourceStatusResponse struct {
	// in: body
	Body map[string]uint8
}

// swagger:response removeMemberSuccessResponse
type removeMemberSuccessResponse struct {
	Body endpointSuccess
}

// swagger:response endpointSuccessResponse
type endpointSuccess struct {
	Ack bool `json:"ack"`
}

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
	j, _ := json.Marshal(endpointSuccess{
		Ack: true,
	})
	w.Write(j)
}

func (rs resourceAPI) addMultipleMembersToResource(w http.ResponseWriter, req *http.Request) {
	var membersResource membersResourceRelation

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
	var update memberResourceRelation

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
	j, _ := json.Marshal(endpointSuccess{
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
	j, _ := json.Marshal(endpointSuccess{
		Ack: true,
	})
	w.Write(j)
}

func (rs resourceAPI) deleteResourceACL(w http.ResponseWriter, req *http.Request) {
	resourcemanager.DeleteResourceACL()

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(endpointSuccess{
		Ack: true,
	})
	w.Write(j)
}
