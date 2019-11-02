package postgres

import (
	"time"

	core "bitbucket.org/picolotec/realtime-location-integrator/internal/core"
	"bitbucket.org/picolotec/realtime-location-integrator/internal/out"
	"github.com/lib/pq"
)

const NOTIFY_PG = "realtime_location_record"

func errorReporter(ev pq.ListenerEventType, err error) {
	out.FailOnError(err, "Failed to create the listener")
}

func Listener(config core.Config, notify chan string) {
	listener := pq.NewListener(config.Server.Postgres.Url, 10*time.Second, time.Minute, errorReporter)
	err := listener.Listen(NOTIFY_PG)
	out.FailOnError(err, "Failed to run listener")

	for {
		select {
		case notification := <-listener.Notify:
			notify <- notification.Extra

		case <-time.After(90 * time.Second):
			go func() {
				err := listener.Ping()
				out.FailOnError(err, "Failed to ping")
			}()
		}
	}
}
