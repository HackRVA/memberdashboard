package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"memberserver/internal/controllers"
	"memberserver/internal/controllers/auth"
	"memberserver/internal/datastore"
	"memberserver/internal/datastore/dbstore"
	"memberserver/internal/datastore/in_memory"
	router "memberserver/internal/routes"
	"memberserver/internal/services/config"
	"memberserver/internal/services/logger"
	"memberserver/internal/services/scheduler"
	"memberserver/internal/services/scheduler/jobs"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

// setupDB
// if the configuration does not include a DBConnectionString
//
//	the server will run with a fakeDB that runs in memory
func setupDB() (datastore.DataStore, error) {
	c, _ := config.Load()
	if c.DBConnectionString == "" {
		return in_memory.Setup()
	}
	return dbstore.Setup()
}

func main() {
	db, err := setupDB()
	if err != nil {
		log.Fatal("error setting up db: %s", err)
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
