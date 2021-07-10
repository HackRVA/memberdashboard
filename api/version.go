package api

import (
	"encoding/json"
	"memberserver/api/models"
	"net/http"
)

// GitCommit is populated by a golang build arg
var GitCommit string

func (a API) getVersion(w http.ResponseWriter, req *http.Request) {
	version := &models.VersionResponse{
		Major:  "0",
		Minor:  "0",
		Hotfix: "0",
		Build:  GitCommit,
	}

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(version)
	w.Write(j)
}
