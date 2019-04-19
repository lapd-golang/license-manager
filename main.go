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

	flag.Parse()

	// Attempt to connect to rabbitmq upon start up
	conn, err := rabbit.Connect(*amqpURI)
	if err != nil {
		log.Fatal(err)
	}

	r := conn.Consume()

	// Attempt to connect to the sqllite3 database
	dbc, err := db.Connect("user_licenses.db")
	defer dbc.DB.Close()

	go conn.Handle(r, dbc)
	// go func(deliveries <-chan amqp.Delivery) {
	// 	for d := range deliveries {
	// 		log.Printf(
	// 			"got %dB delivery: [%v] %s",
	// 			len(d.Body),
	// 			d.DeliveryTag,
	// 			d.Body,
	// 		)
	// 	}
	// }(r)

	//TODO go routine which listens to rabbitMQ messages

	router, err := Routes(dbc)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Now listening on localhost:8080..")
	log.Fatal(http.ListenAndServe(":8080", router))

}
