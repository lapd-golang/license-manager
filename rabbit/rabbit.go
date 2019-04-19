package rabbit

import (
	"encoding/json"
	"fmt"

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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func Connect(amqpURI string) (*RMQConn, error) {
	// Attempt to connect to RabbitMQ, Hopefully it is running
	// TODO: Create a retry ticker for the inital connection.
	conn, err := amqp.Dial(amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Infof("Connected to RabbitMQ on  %s", amqpURI)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		EXCHANGE, // name
		"topic",  // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // noWait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare the Exchange")
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
	failOnError(err, "Error declaring the Queue")

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		q.Name, q.Messages, q.Consumers, "go-test-key")

	err = ch.QueueBind(
		q.Name,   // name of the queue
		"a-key",  // bindingKey
		EXCHANGE, // sourceExchange
		false,    // noWait
		nil,      // arguments
	)
	failOnError(err, "Error binding to the Queue")

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", "go-amqp-example")

	return &RMQConn{ch, q, EXCHANGE}, nil
}

func (c *RMQConn) Consume() <-chan amqp.Delivery {
	var replies <-chan amqp.Delivery
	replies, err := c.ch.Consume(
		c.q.Name,          // queue
		"go-amqp-example", // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	failOnError(err, "Error consuming the Queue")
	return replies
}

//Handle -  These messages are autoacked.
// Recieves the messages and then inserts into the database
func (c *RMQConn) Handle(deliveries <-chan amqp.Delivery, dbc *db.Dao) {

	for d := range deliveries {
		m := Msg{}
		json.Unmarshal(d.Body, &m)
		log.Infof("%s", m.Code)
		log.Printf(
			"got %dB delivery: [%v] %s",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
	}
}
