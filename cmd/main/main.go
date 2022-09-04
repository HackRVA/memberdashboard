package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"memberserver/internal/controllers"
	"memberserver/internal/controllers/auth"
	"memberserver/internal/datastore/dbstore"
	router "memberserver/internal/routes"
	"memberserver/internal/services/logger"
	"memberserver/internal/services/scheduler"
	"memberserver/internal/services/scheduler/jobs"
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
	api := controllers.Setup(db, auth)
	router := router.New(api, auth)

	srv := &http.Server{
		Handler: router.UnAuthedRouter,
		Addr:    "0.0.0.0:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	j := jobs.New(db, logger.New())
	s := scheduler.Scheduler{}

	go s.Setup(j)

	log.Debug("Server listening on http://localhost:3000/")
	log.Fatal(srv.ListenAndServe())
}
