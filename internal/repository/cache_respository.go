package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	redisClient *redis.Client
}

const (
	EXPIRE_TIME = time.Hour * 8
)

var CTX = context.Background()

func NewCacheRepository(redisClient *redis.Client) *CacheRepository {
	return &CacheRepository{
		redisClient: redisClient,
	}
}

func (repo *CacheRepository) Set(key string, val string) error {
	err := repo.redisClient.Set(CTX, key, val, EXPIRE_TIME).Err()
	return err
}

func (repo *CacheRepository) Get(key string) (string, error) {
	val, err := repo.redisClient.Get(CTX, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
