package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	core "bitbucket.org/picolotec/realtime-location-integrator/internal/core"
	"bitbucket.org/picolotec/realtime-location-integrator/internal/out"
	"github.com/go-redis/redis"
	"github.com/labstack/gommon/color"
)

type RealtimeLocation struct {
	License string `json:"license_number"`
}

func Connect(config core.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Server.Redis.Url,
		Password: config.Server.Redis.Password,
		DB:       0,
	})
	_, err := client.Ping().Result()
	out.FailOnError(err, "Failed to connect to Redis")

	return client
}

func Sent(config core.Config, payload string, client *redis.Client) {
	key := redisKey(payload, config.Client)

	err := client.Set(key, payload, config.Server.Expiration).Err()
	out.FailOnError(err, "Failed to store in redis: "+config.Client)

	log.Printf("[%s] - %s", color.Green("redis"), config.Client)
}

func redisKey(data string, client string) string {
	var rl RealtimeLocation

	json.Unmarshal([]byte(data), &rl)
	sec := time.Now().Unix()

	return fmt.Sprintf("%s:%s:%d", client, rl.License, sec)
}
