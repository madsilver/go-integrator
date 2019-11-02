package core

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"bitbucket.org/picolotec/realtime-location-integrator/internal/out"
)

type Server struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Url      string
}

type Config struct {
	Server struct {
		Postgres   Server        `yaml:"postgres"`
		Rabbitmq   Server        `yaml:"rabbitmq"`
		Redis      Server        `yaml:"redis"`
		Expiration time.Duration `yaml:"expiration"`
	}
	Client string
	Events struct {
		Operational []string
		System      []string
	}
}

func Load(client string, env string) Config {
	config := Config{}
	config.Client = client

	configFile := "config/config.yml"
	if env == "dev" {
		configFile = "config/config.dev.yml"
	}

	ev, err := os.Open(configFile)
	out.FailOnError(err, "")

	decoder := yaml.NewDecoder(ev)
	err = decoder.Decode(&config)

	config.Server.Postgres.Url = fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Server.Postgres.Username,
		config.Server.Postgres.Password,
		config.Server.Postgres.Host,
		config.Server.Postgres.Port,
		client,
	)

	config.Server.Rabbitmq.Url = fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		config.Server.Rabbitmq.Username,
		config.Server.Rabbitmq.Password,
		config.Server.Rabbitmq.Host,
		config.Server.Rabbitmq.Port,
	)

	config.Server.Expiration = time.Duration(config.Server.Expiration) * time.Second
	config.Server.Redis.Url = fmt.Sprintf(
		"%s:%d",
		config.Server.Redis.Host,
		config.Server.Redis.Port,
	)

	return config
}
