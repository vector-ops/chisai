package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() error {
	uri := os.Getenv("REDIS_URL")
	if uri == "" {
		return errors.New("REDIS_URL env var not set")
	}

	opts, err := redis.ParseURL(uri)
	if err != nil {
		return err
	}

	client := redis.NewClient(opts)

	if err = client.Ping(context.Background()).Err(); err != nil {
		return err
	}

	rdb = client

	return nil
}

func GetRedis() *redis.Client {
	return rdb
}

func CloseRedis(ctx context.Context) error {
	return rdb.Close()
}
