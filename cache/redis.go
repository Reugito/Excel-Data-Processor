package cache

import (
	"context"
	"dataProcessor/config"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisRepo struct {
	Client *redis.Client
}

// ConnectRedis connects to Redis cache
func ConnectRedis(cfg config.RedisConfig) (*RedisRepo, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	fmt.Println(pong, err)

	fmt.Println("Connected to Redis")

	return &RedisRepo{Client: rdb}, nil
}
