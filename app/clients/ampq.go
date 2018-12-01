package clients

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Ampq struct {
	uri string
}

var queueName = "tt3398228"

type RequestMessage struct {
	Id        int64  `json:"id" binding:"required"`
	Message   string `json:"message" binding:"required"`
	CreatedAt int64  `json:"created" binding:"required"`
}

func (omdb Omdb) SendToQueue(id int64, created string, message string) {
	var err error
	var conn = &amqp.Connection{}

	if conn, err = amqp.Dial(omdb.uri); err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()
	var ch = &amqp.Channel{}
	if ch, err = conn.Channel(); err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare to queue: %s", err)
	}

	layout := "Mon Jan 02 15:04:05 -0700 2006"
	date := created
	t, err := time.Parse(layout, date)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t.Unix())

	tweet := RequestMessage{
		Id:        id,
		CreatedAt: t.Unix(),
		Message:   message,
	}
	res1B, _ := json.Marshal(tweet)
	body := string(res1B)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}
}
