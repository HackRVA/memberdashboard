package api

import (
	"encoding/json"
	"memberserver/api/models"
	"net/http"
)

// GitCommit is populated by a golang build arg
var GitCommit string

func (a API) getVersion(w http.ResponseWriter, req *http.Request) {
	var version models.VersionResponse

	version.Commit = GitCommit

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(version)
	w.Write(j)
}
