package v1

import (
	"encoding/json"
	"net/http"

	"github.com/HackRVA/memberserver/models"
	"github.com/sirupsen/logrus"
)

// GitCommit is populated by a golang build arg
var GitCommit string

// NewInMemoryVersionStore initialises an empty version store.
func NewInMemoryVersionStore() *InMemoryVersionStore {
	return &InMemoryVersionStore{}
}

type InMemoryVersionStore struct{}

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
	switch r.Method {
	case http.MethodGet:
		v.showVersion(w)
	}
}

func (v *VersionServer) showVersion(w http.ResponseWriter) {
	var versionInfo models.VersionResponse

	err := json.Unmarshal(v.store.GetVersion(), &versionInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(versionInfo.Major) == 0 || len(versionInfo.Build) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/text")
		if _, err := w.Write([]byte("some issue getting the version")); err != nil {
			logrus.Error(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(v.store.GetVersion()); err != nil {
		logrus.Error(err)
	}
}
