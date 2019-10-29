package main

import (
	"log"
)

func main() {
	config := loadEnv()
	log.Print("Listener running...")
	log.Printf("Client: %s", config.client)

	run(config)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
