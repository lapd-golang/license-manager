package main

import (
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sevren/test/db"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Attempt to connect to the sqllite3 database
	dbc, err := db.Connect("user_licenses.db")
	defer dbc.DB.Close()

	//TODO go routine which listens to rabbitMQ messages

	router, err := Routes(dbc)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Now listening on localhost:8080..")
	log.Fatal(http.ListenAndServe(":8080", router))

}
