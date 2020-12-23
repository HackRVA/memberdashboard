package main

import (
	"log"
	"net/http"

	"github.com/dfirebaugh/memberserver/config"
	"github.com/dfirebaugh/memberserver/routes"
)

func main() {
	routes.Setup()
	c, err := config.Load("./sample.config.json")

	if err != nil {
		log.Fatal(err)
	}

	log.Println(c.TestValue)
	log.Println(c.SomethingElse)

	log.Print("Server listening on http://localhost:3000/")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
