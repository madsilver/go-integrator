package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config contains various config data populated from YAML
type Config struct {
	PostgresURL string
	RabbitMQURL string
}

func main() {
	config := Config{}
	filebytes, err := ioutil.ReadFile("prcreds.yaml")
	if err != nil {
		log.Fatal("Failed to read creds")
	}

	err = yaml.Unmarshal(filebytes, &config)
	if err != nil {
		log.Fatal("Failed to parse creds", err)
	}
}
