package clients

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)

type Ampq struct {
	uri        string
	connection *amqp.Connection
	channel    *amqp.Channel
}

var (
	ampqInstance *Ampq
	initAmpqOnce sync.Once
)

func (m Ampq) Get() *Ampq {
	initAmpqOnce.Do(func() {
		var err error
		if m.connection, err = amqp.Dial(m.uri); err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		}
		//defer conn.Close()
		if m.channel, err = m.connection.Channel(); err != nil {
			log.Fatalf("Failed to open a channel: %s", err)
		}
		ampqInstance = &m
		//defer ch.Close()
	})
	return ampqInstance
}

type RequestMessage struct {
	Id        int64  `json:"id" binding:"required"`
	Message   string `json:"message" binding:"required"`
	CreatedAt int64  `json:"created" binding:"required"`
}

func (c Ampq) PushMessage(group string, id int64, created string, message string) (err error) {
	q, err := c.channel.QueueDeclare(
		group, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare to queue: %s", err)
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
	err = c.channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %s", err)
	}
	return
}

type TaskMessage struct {
	Type int `json:"id" binding:"required"`
}

func (c Ampq) PushTask(group string, taskType int) (err error) {
	q, err := c.channel.QueueDeclare(
		group, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare to queue: %s", err)
	}

	task := TaskMessage{
		Type: taskType,
	}
	res1B, _ := json.Marshal(task)
	body := string(res1B)
	err = c.channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %s", err)
	}
	return
}
