package rabbitmq

import (
	"log"

	core "bitbucket.org/picolotec/realtime-location-integrator/internal/core"
	"bitbucket.org/picolotec/realtime-location-integrator/internal/out"
	"github.com/labstack/gommon/color"
	"github.com/streadway/amqp"
	"github.com/tidwall/sjson"
)

func Connect(conf core.Config) *amqp.Channel {
	config := conf

	conn, err := amqp.Dial(config.Server.Rabbitmq.Url)
	out.FailOnError(err, "Failed to connect to RabbitMQ")

	channel, err := conn.Channel()
	out.FailOnError(err, "Failed to open a channel")

	return channel
}

func Publish(payload string, channel *amqp.Channel, event string, client string) {
	queue, err := channel.QueueDeclare(
		event,
		false,
		false,
		false,
		false,
		nil,
	)
	out.FailOnError(err, "Failed to declare a queue")

	body := clientPayload(payload, client)

	err2 := channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	out.FailOnError(err2, "Failed to publish a message: "+event)

	log.Printf("[%s] - %s - %s", color.Green("rabbitmq"), client, event)
}

func clientPayload(payload string, client string) string {

	body, _ := sjson.Set(payload, "client", client)

	return body
}
