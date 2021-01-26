package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"memberserver/api"
	"memberserver/database"
	"memberserver/resourcemanager"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}

	rm, err := resourcemanager.Setup()
	if err != nil {
		log.Errorf("error setting up resource manager: %s", err)
	}

	router := api.Setup(db, rm)

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Debug("Server listening on http://localhost:3000/")
	log.Fatal(srv.ListenAndServe())
}
