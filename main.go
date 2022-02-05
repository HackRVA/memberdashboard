package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"memberserver/api"
	"memberserver/api/auth"
	"memberserver/api/router"
	"memberserver/datastore/dbstore"
	"memberserver/scheduler"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	db, err := dbstore.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}

	auth := auth.New(db)
	api := api.Setup(db, auth)
	router := router.New(api, auth)

	srv := &http.Server{
		Handler: router.UnAuthedRouter,
		Addr:    "0.0.0.0:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	schedule := scheduler.Scheduler{}

	go schedule.Setup(db)

	log.Debug("Server listening on http://localhost:3000/")
	log.Fatal(srv.ListenAndServe())
}
