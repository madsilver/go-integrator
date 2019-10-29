package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	pq "github.com/lib/pq"
	"github.com/streadway/amqp"
)

type Config struct {
	postgres string
	rabbitmq string
}

func init() {
	err := godotenv.Load()
	failOnError(err, "No .env file found")
}

func main() {
	config := loadEnv()
	log.Print("Listener running...")
	run(config)
}

func errorReporter(ev pq.ListenerEventType, err error) {
	failOnError(err, "Failed to create the listener")
}

func run(config Config) {
	listener := pq.NewListener(config.postgres, 10*time.Second, time.Minute, errorReporter)
	err := listener.Listen("realtime_location_record")
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
			"cidadedutra.realtime_location",
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

func loadEnv() Config {
	config := Config{}

	config.postgres = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("POSTGRES_USERNAME"),
		getEnv("POSTGRES_PASSWORD"),
		getEnv("POSTGRES_HOST"),
		getEnv("POSTGRES_PORT"),
		"cidadedutra",
	)

	config.rabbitmq = fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		getEnv("RABBITMQ_USERNAME"),
		getEnv("RABBITMQ_PASSWORD"),
		getEnv("RABBITMQ_HOST"),
		getEnv("RABBITMQ_PORT"),
	)

	return config
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return ""
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
