package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dfirebaugh/memberserver/config"
	"github.com/dfirebaugh/memberserver/routes"
)

func main() {
	routes.Setup()

	c, err := config.Load(os.Getenv("MEMBER_SERVER_CONFIG_FILE"))

	if err != nil {
		log.Fatal(err)
	}

	log.Println(c.PaypalURL)
	log.Println(c.PaypalSignature)

	log.Print("Server listening on http://localhost:3000/")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
