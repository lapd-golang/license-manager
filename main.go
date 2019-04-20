package main

import (
	"flag"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sevren/test/db"
	"github.com/sevren/test/rabbit"
	log "github.com/sirupsen/logrus"
)

var (
	amqpURI = flag.String("amqp", "amqp://guest:guest@localhost:5672/", "AMQP URI")
)

func main() {

	challenge3features := false
	flag.Parse()

	// Attempt to connect to the sqllite3 database
	dbc, err := db.Connect("user_licenses.db")
	defer dbc.DB.Close()

	// Attempt to connect to rabbitmq upon start up
	// If rabbitmq can not be connected then challenge 3 stuff is disabled
	// You can still use the REST interface
	conn, err := rabbit.Connect(*amqpURI)
	if err != nil {
		log.Warnf("Challenge 3 features disabled. Could not connect to rabbitmq \n", err)
	}

	if conn != nil {
		r := conn.Consume()
		go conn.Handle(r, dbc)
		challenge3features = true
	}

	if challenge3features {
		log.Info("Challenge 3 features enabled. Licenses will generate more uniquely")
	}

	router, err := Routes(dbc, challenge3features)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Now listening on localhost:8080..")
	log.Fatal(http.ListenAndServe(":8080", router))

}
