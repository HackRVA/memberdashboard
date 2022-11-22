package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/HackRVA/memberserver/internal/models"
)

// GitCommit is populated by a golang build arg
var GitCommit string

// NewInMemoryVersionStore initialises an empty version store.
func NewInMemoryVersionStore() *InMemoryVersionStore {
	return &InMemoryVersionStore{}
}

type InMemoryVersionStore struct {
}

func (i *InMemoryVersionStore) GetVersion() []byte {
	if len(GitCommit) == 0 {
		GitCommit = "dev"
	}
	version := models.VersionResponse{
		Major: "1",
		Build: GitCommit,
	}
	j, _ := json.Marshal(version)

	return j
}

// VersionStore stores version information about the app.
type VersionStore interface {
	GetVersion() []byte
}

// VersionServer is a HTTP interface for version information.
type VersionServer struct {
	store VersionStore
}

func (v *VersionServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	version := strings.TrimPrefix(r.URL.Path, "version")
	switch r.Method {
	case http.MethodGet:
		v.showVersion(w, version)
	}
}

func (v *VersionServer) showVersion(w http.ResponseWriter, version string) {
	var versionInfo models.VersionResponse

	err := json.Unmarshal(v.store.GetVersion(), &versionInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(versionInfo.Major) == 0 || len(versionInfo.Build) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("some issue getting the version"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(v.store.GetVersion())
}
