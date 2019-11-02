package main

import (
	"flag"
	"log"

	"bitbucket.org/picolotec/realtime-location-integrator/internal/core"
	"github.com/labstack/gommon/color"
)

func main() {
	client := flag.String("client", "", "The name of the client database")
	env := flag.String("env", "prod", "The name of the environment")
	flag.Parse()

	var config core.Config
	config = core.Load(*client, *env)

	log.Print("Integrator running...")
	log.Printf("Client: %s", color.Green(config.Client))

	run(config)
}
