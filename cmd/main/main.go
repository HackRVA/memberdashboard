package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	router "github.com/HackRVA/memberserver/internal/routes"
	"github.com/HackRVA/memberserver/pkg/mqtt"
	"github.com/HackRVA/memberserver/pkg/paypal"
	"github.com/HackRVA/memberserver/pkg/slack"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/internal/controllers"
	"github.com/HackRVA/memberserver/internal/controllers/auth"
	"github.com/HackRVA/memberserver/internal/datastore"
	"github.com/HackRVA/memberserver/internal/datastore/dbstore"
	"github.com/HackRVA/memberserver/internal/datastore/in_memory"
	"github.com/HackRVA/memberserver/internal/services/logger"
	"github.com/HackRVA/memberserver/internal/services/member"
	"github.com/HackRVA/memberserver/internal/services/resourcemanager"
	"github.com/HackRVA/memberserver/internal/services/scheduler"
	"github.com/HackRVA/memberserver/internal/services/scheduler/jobs"
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
		log.Fatalf("error setting up db: %s", err)
	}
	c := config.Get()

	log := logger.New()
	rm := resourcemanager.New(mqtt.New(), db, slack.Notifier{WebHookURL: c.SlackAccessEvents}, log)
	pp := paypal.Setup(c.PaypalURL, c.PaypalClientID, c.PaypalClientSecret, log)

	auth := auth.New(db)
	api := controllers.Setup(db, auth, rm, pp, log)
	router := router.New(api, auth)

	srv := &http.Server{
		Handler: router.UnAuthedRouter,
		Addr:    "0.0.0.0:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	j := jobs.New(db, log, member.New(db, rm, pp, log), rm)
	s := scheduler.Scheduler{}

	go s.Setup(j)

	log.Debug("Server listening on http://localhost:3000/")
	log.Fatal(srv.ListenAndServe())
}
