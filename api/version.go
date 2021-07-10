package api

import (
	"encoding/json"
	"memberserver/api/models"
	"net/http"
)

func (a API) getVersion(w http.ResponseWriter, req *http.Request) {
	var version models.VersionResponse

	version.Commit = a.gitCommit

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(version)
	w.Write(j)
}
