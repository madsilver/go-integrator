package main

import (
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/streadway/amqp"
)

const NOTIFY_PG = "realtime_location_record"

func errorReporter(ev pq.ListenerEventType, err error) {
	failOnError(err, "Failed to create the listener")
}

func run(config Config) {
	listener := pq.NewListener(config.postgres, 10*time.Second, time.Minute, errorReporter)
	err := listener.Listen(NOTIFY_PG)
	failOnError(err, "Failed to run listener")

	rabbitchannel := make(chan string, 100)

	go func() {
		conn, err := amqp.Dial(config.rabbitmq)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			config.clientQueue,
			false,
			false,
			false,
			false,
			nil,
		)
		failOnError(err, "Failed to declare a queue")

		for {
			payload := <-rabbitchannel
			err := ch.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(payload),
				})
			log.Printf(" [x] Sent %s", payload)
			failOnError(err, "Failed to publish a message")
		}
	}()

	for {
		select {
		case notification := <-listener.Notify:
			rabbitchannel <- notification.Extra

		case <-time.After(90 * time.Second):
			go func() {
				err := listener.Ping()
				failOnError(err, "Failed to ping")
			}()
		}
	}
}
