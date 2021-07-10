package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"memberserver/api"
	"memberserver/database"
	"memberserver/scheduler"
)

// GitCommit must be in the main go file in order to get the commit hash
var GitCommit string

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}
	defer db.Release()

	router := api.Setup(db, GitCommit)

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go scheduler.Setup(db)

	log.Debug("Server listening on http://localhost:3000/")
	log.Fatal(srv.ListenAndServe())
}
