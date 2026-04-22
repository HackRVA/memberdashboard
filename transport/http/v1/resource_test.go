package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HackRVA/memberserver/datastore/dbstore"
	"github.com/HackRVA/memberserver/datastore/in_memory"
	"github.com/HackRVA/memberserver/models"

	"github.com/gorilla/mux"
)

// resourceNameStore embeds in_memory but implements GetResourceByName by
// returning a canned value -- in_memory's stub returns an empty Resource{}
// which defeats heartbeat lookups keyed by name.
type resourceNameStore struct {
	*in_memory.In_memory
	resource models.Resource
	err      error
}

func (s *resourceNameStore) GetResourceByName(ctx context.Context, name string) (models.Resource, error) {
	if s.err != nil {
		return models.Resource{}, s.err
	}
	return s.resource, nil
}

func TestResourceStatus_MissingName(t *testing.T) {
	api := resourceAPI{db: &in_memory.In_memory{}}

	req := httptest.NewRequest(http.MethodGet, "/api/resource/status/", nil)
	rec := httptest.NewRecorder()

	api.ResourceStatus(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestResourceStatus_NotFound(t *testing.T) {
	api := resourceAPI{db: &resourceNameStore{
		In_memory: &in_memory.In_memory{},
		err:       errNotFoundStub,
	}}

	req := httptest.NewRequest(http.MethodGet, "/api/resource/status/ghost", nil)
	req = mux.SetURLVars(req, map[string]string{"name": "ghost"})
	rec := httptest.NewRecorder()

	api.ResourceStatus(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestResourceStatus_Offline(t *testing.T) {
	api := resourceAPI{db: &resourceNameStore{
		In_memory: &in_memory.In_memory{},
		resource:  models.Resource{Name: "frontdoor-offline"},
	}}

	req := httptest.NewRequest(http.MethodGet, "/api/resource/status/frontdoor-offline", nil)
	req = mux.SetURLVars(req, map[string]string{"name": "frontdoor-offline"})
	rec := httptest.NewRecorder()

	api.ResourceStatus(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
	if rec.Body.String() != "offline" {
		t.Errorf("body = %q, want %q", rec.Body.String(), "offline")
	}
}

func TestResourceStatus_Ok(t *testing.T) {
	api := resourceAPI{db: &resourceNameStore{
		In_memory: &in_memory.In_memory{},
		resource:  models.Resource{Name: "frontdoor-fresh"},
	}}

	dbstore.ResourceHeartbeat(models.Resource{Name: "frontdoor-fresh"})

	req := httptest.NewRequest(http.MethodGet, "/api/resource/status/frontdoor-fresh", nil)
	req = mux.SetURLVars(req, map[string]string{"name": "frontdoor-fresh"})
	rec := httptest.NewRecorder()

	api.ResourceStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "ok" {
		t.Errorf("body = %q, want %q", rec.Body.String(), "ok")
	}
}

type stubError string

func (s stubError) Error() string { return string(s) }

const errNotFoundStub = stubError("not found")
