package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/HackRVA/memberserver/pkg/membermgr/models"
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
	resources := rs.db.GetResources()
	ok(w, resources)
}

func (rs resourceAPI) update(w http.ResponseWriter, req *http.Request) {
	var updateResourceReq models.Resource

	err := json.NewDecoder(req.Body).Decode(&updateResourceReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r, err := rs.db.UpdateResource(updateResourceReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok(w, r)
}

func (rs resourceAPI) delete(w http.ResponseWriter, req *http.Request) {
	var deleteResourceReq models.ResourceDeleteRequest

	err := json.NewDecoder(req.Body).Decode(&deleteResourceReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rs.logger.Printf("attempting to delete %s", deleteResourceReq.ID)

	err = rs.db.DeleteResource(deleteResourceReq.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok(w, models.EndpointSuccess{
		Ack: true,
	})
}

func (rs resourceAPI) AddMultipleMembersToResource(w http.ResponseWriter, req *http.Request) {
	var membersResource models.MembersResourceRelation

	err := json.NewDecoder(req.Body).Decode(&membersResource)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resource, err := rs.db.AddMultipleMembersToResource(membersResource.Emails, membersResource.ID)
	for _, email := range membersResource.Emails {
		member, _ := rs.db.GetMemberByEmail(email)
		rs.logger.Info("pushing member to resource", member.Email, member.Resources)
		rs.resourcemanager.PushOne(member)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok(w, resource)
}

func (rs resourceAPI) RemoveMember(w http.ResponseWriter, req *http.Request) {
	var update models.MemberResourceRelationUpdateRequest

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

	ok(w, models.EndpointSuccess{
		Ack: true,
	})

	resource, err := rs.db.GetResourceByID(update.ID)
	if err != nil {
		rs.logger.Errorf("error getting resource to update when removing a member: %s", err)
	}

	rs.resourcemanager.UpdateResourceACL(resource)
	rs.resourcemanager.UpdateResources()
}

func (rs resourceAPI) Register(w http.ResponseWriter, req *http.Request) {
	var register models.RegisterResourceRequest

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

	ok(w, r)
}

func (rs resourceAPI) Status(w http.ResponseWriter, req *http.Request) {
	resources := rs.db.GetResources()
	// statusMap := make(map[string]uint8)

	for _, r := range resources {
		if r == (models.Resource{}) {
			continue
		}
		rs.resourcemanager.CheckStatus(r)
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

func (rs resourceAPI) UpdateResourceACL(w http.ResponseWriter, req *http.Request) {
	rs.resourcemanager.UpdateResources()

	ok(w, models.EndpointSuccess{
		Ack: true,
	})
}

func (rs resourceAPI) Open(w http.ResponseWriter, req *http.Request) {
	var openResourceRequest models.OpenResourceRequest

	err := json.NewDecoder(req.Body).Decode(&openResourceRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resource, err := rs.db.GetResourceByName(openResourceRequest.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rs.resourcemanager.Open(resource)

	ok(w, models.EndpointSuccess{
		Ack: true,
	})
}

func (rs resourceAPI) DeleteResourceACL(w http.ResponseWriter, req *http.Request) {
	rs.resourcemanager.DeleteResourceACL()

	ok(w, models.EndpointSuccess{
		Ack: true,
	})
}
