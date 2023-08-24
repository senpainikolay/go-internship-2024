package redisdb

import (
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {

	intVar, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Could not upload db_nr from .env")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DB:       intVar,
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	log.Printf("successfully connect to redis db...\n")
	return client
}
