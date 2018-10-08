package clients

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

var queueName = "tt3398228"

type RequestMessage struct {
	Id        int64  `json:"id" binding:"required"`
	Message   string `json:"message" binding:"required"`
	CreatedAt int64  `json:"created" binding:"required"`
}

func SendToQueue(id int64, created string, message string) {
	conn, err := amqp.Dial(config.AmpqUri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

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
	failOnError(err, "Failed to publish a message")
}
