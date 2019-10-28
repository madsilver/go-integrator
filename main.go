package main

import (
	"io/ioutil"
	"log"
	"time"

	pq "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

// Config contains various config data populated from YAML
type Config struct {
	PostgresURL string
	RabbitMQURL string
}

func main() {
	config := Config{}
	filebytes, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Failed to read creds")
	}

	err = yaml.Unmarshal(filebytes, &config)
	if err != nil {
		log.Fatal("Failed to parse creds", err)
	}

	run(config)
}

func errorReporter(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Print(err)
	}
}

func run(config Config) {
	listener := pq.NewListener(config.PostgresURL, 10*time.Second, time.Minute, errorReporter)
	err := listener.Listen("realtime_location_record")
	if err != nil {
		log.Fatal(err)
	}

}
