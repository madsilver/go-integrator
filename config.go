package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	postgres        string
	postgresTrigger string
	rabbitmq        string
	client          string
	clientQueue     string
}

func init() {
	err := godotenv.Load()
	failOnError(err, "No .env file found")
}

func loadEnv() Config {
	config := Config{}

	client := flag.String("client", "null", "a string")
	flag.Parse()

	config.postgres = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("POSTGRES_USERNAME"),
		getEnv("POSTGRES_PASSWORD"),
		getEnv("POSTGRES_HOST"),
		getEnv("POSTGRES_PORT"),
		*client,
	)

	config.rabbitmq = fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		getEnv("RABBITMQ_USERNAME"),
		getEnv("RABBITMQ_PASSWORD"),
		getEnv("RABBITMQ_HOST"),
		getEnv("RABBITMQ_PORT"),
	)

	config.client = *client
	config.clientQueue = fmt.Sprintf("%s.realtime_location", *client)

	return config
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return ""
}
