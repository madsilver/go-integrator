package main

import (
	"fmt"

	core "bitbucket.org/picolotec/realtime-location-integrator/internal/core"
	"bitbucket.org/picolotec/realtime-location-integrator/internal/postgres"
	"bitbucket.org/picolotec/realtime-location-integrator/internal/rabbitmq"
	"bitbucket.org/picolotec/realtime-location-integrator/internal/redis"
)

const EVENT_PREFIX = "analyze"

func run(config core.Config) {
	notify := make(chan string, 100)

	go postgres.Listener(config, notify)

	ch := rabbitmq.Connect(config)
	cl := redis.Connect(config)

	for {
		payload := <-notify

		go redis.Sent(config, payload, cl)

		for _, op := range config.Events.Operational {
			eventOp := fmt.Sprintf("%s.%s", EVENT_PREFIX, op)
			rabbitmq.Publish(payload, ch, eventOp, config.Client)
		}

		for _, sys := range config.Events.System {
			eventSys := fmt.Sprintf("%s.%s", EVENT_PREFIX, sys)
			rabbitmq.Publish(payload, ch, eventSys, config.Client)
		}

	}
}
