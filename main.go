package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/HackRVA/memberserver/pkg/mqtt"
	"github.com/HackRVA/memberserver/pkg/paypal"
	"github.com/HackRVA/memberserver/pkg/slack"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/datastore"
	"github.com/HackRVA/memberserver/datastore/dbstore"
	"github.com/HackRVA/memberserver/datastore/in_memory"
	"github.com/HackRVA/memberserver/services/logger"
	"github.com/HackRVA/memberserver/services/member"
	"github.com/HackRVA/memberserver/services/scheduler"
	"github.com/HackRVA/memberserver/services/scheduler/jobs"
	httprouter "github.com/HackRVA/memberserver/transport/http"
	"github.com/HackRVA/memberserver/transport/http/middleware/auth"
	v1 "github.com/HackRVA/memberserver/transport/http/v1"
	mqtthandler "github.com/HackRVA/memberserver/transport/mqtt/v1"
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
	mqttHandler := mqtthandler.New(mqtt.New(), db, slack.Notifier{WebHookURL: c.SlackAccessEvents})
	pp := paypal.Setup(c.PaypalURL, c.PaypalClientID, c.PaypalClientSecret, log)

	auth := auth.New(db)
	muxRouter := mux.NewRouter()
	v1Router := v1.New(muxRouter, db, auth, mqttHandler, pp, log)
	router := httprouter.New(muxRouter, []httprouter.MemberServerRouter{v1Router})

	srv := &http.Server{
		Handler: router.UnAuthedRouter,
		Addr:    "0.0.0.0:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	j := jobs.New(db, log, member.New(db, mqttHandler, pp), mqttHandler)
	s := scheduler.Scheduler{}

	go s.Setup(j)

	log.Debug("Server listening on http://localhost:3000/")
	log.Fatal(srv.ListenAndServe())
}
