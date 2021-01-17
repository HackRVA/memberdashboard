package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dfirebaugh/memberserver/config"
	"github.com/dfirebaugh/memberserver/routes"
)

func main() {
	router := routes.Setup()

	println(os.Getenv("MEMBER_SERVER_CONFIG_FILE"))

	_, err := config.Load(os.Getenv("MEMBER_SERVER_CONFIG_FILE"))
	if len(os.Getenv("MEMBER_SERVER_CONFIG_FILE")) == 0 {
		log.Fatal("must set the MEMBER_SERVER_CONFIG_FILE environment variable to point to config file")
	}

	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Server listening on http://localhost:3000/")
	log.Fatal(srv.ListenAndServe())
}
