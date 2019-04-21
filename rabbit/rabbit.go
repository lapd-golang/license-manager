package rabbit

//This file is responsible for providing all the functionality for connecting and consuming messages from RabbitMQ
// RabbitMQ is part of the Challenge 3.

import (
	"encoding/json"

	"github.com/sevren/test/db"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const EXCHANGE = "data"

type RMQConn struct {
	ch *amqp.Channel
	q  amqp.Queue
	Ex string
}

type Msg struct {
	Code string `json:"code"`
}

func warnOnError(err error, msg string) {
	if err != nil {
		log.Warnf("%s: %s", msg, err)
	}
}

func Connect(amqpURI string) (*RMQConn, error) {
	// Attempt to connect to RabbitMQ, Hopefully it is running
	// TODO: Create a retry ticker for the inital connection.
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Warn("Failed to connect to RabbitMQ")
		return nil, err
	}
	log.Infof("Connected to RabbitMQ on  %s", amqpURI)

	ch, err := conn.Channel()
	if err != nil {
		log.Warn("Failed to open a channel")
		return nil, err
	}

	err = ch.ExchangeDeclare(
		EXCHANGE, // name
		"topic",  // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // noWait
		nil,      // arguments
	)
	if err != nil {
		log.Warn("Failed to declare the Exchange")
		return nil, err
	}
	log.Infof("Declaring/Connecting to RabbitMQ Exchange on  %s", EXCHANGE)

	var q amqp.Queue

	q, err = ch.QueueDeclare(
		"go-test-queue", // name, leave empty to generate a unique name
		true,            // durable
		false,           // delete when usused
		false,           // exclusive
		false,           // noWait
		nil,             // arguments
	)
	if err != nil {
		log.Warn("Error declareing the Queue")
		return nil, err
	}

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		q.Name, q.Messages, q.Consumers, "a-key")

	err = ch.QueueBind(
		q.Name,   // name of the queue
		"a-key",  // bindingKey
		EXCHANGE, // sourceExchange
		false,    // noWait
		nil,      // arguments
	)
	if err != nil {
		log.Warn("Error binding to the Queue")
		return nil, err
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", "license-manager")

	return &RMQConn{ch, q, EXCHANGE}, nil
}

func (c *RMQConn) Consume() <-chan amqp.Delivery {
	var replies <-chan amqp.Delivery
	replies, err := c.ch.Consume(
		c.q.Name,          // queue
		"license-manager", // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	warnOnError(err, "Error consuming the Queue")
	return replies
}

//Handle -  These messages are autoacked.
// Recieves the messages and then inserts into the database
func (c *RMQConn) Handle(deliveries <-chan amqp.Delivery, dbc *db.Dao) {

	for d := range deliveries {

		// call the database handler to insert the used code
		m := Msg{}
		json.Unmarshal(d.Body, &m)
		dbc.StoreUsedLicenses(m.Code)

		log.Infof("%s", m.Code)
		log.Printf(
			"got %dB delivery: [%v] %s",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
	}
}
