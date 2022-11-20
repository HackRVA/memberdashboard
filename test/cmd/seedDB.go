package main

import (
	"log"
	"memberserver/internal/datastore/dbstore"
	"memberserver/test/generators"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Provide an argument for # of members to create.")
	}

	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to parse %v as count.", os.Args[1])
	}

	db, _ := dbstore.Setup()
	generators.Seed(db, count)
}
