package cacheRepository

import (
	"context"
	rd "github.com/redis/go-redis/v9"
	"test2/model"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, exp time.Duration) *model.Error
	Get(key string) (string, *model.Error)
	Del(key string) *model.Error
}

type Redis struct {
	context.Context
	*rd.Client
}

func RedisConnection(addr, password string, db int) *rd.Client {
	return rd.NewClient(&rd.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func NewRedis(ctx context.Context, client *rd.Client) Cache {
	return &Redis{
		Context: ctx,
		Client:  client,
	}
}

func (r Redis) Set(key string, value interface{}, exp time.Duration) *model.Error {
	if err := r.Client.Set(r.Context, key, value, exp).Err(); err != nil {
		return model.NewError(500, "Internal Server Error", err)
	}
	return nil
}

func (r Redis) Get(key string) (string, *model.Error) {
	data, err := r.Client.Get(r.Context, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return "", model.NewError(404, "Not Found", err)
		}
		return "", model.NewError(500, "Internal Server Error", err)
	}
	return data, nil
}

func (r Redis) Del(key string) *model.Error {
	if err := r.Client.Del(r.Context, key).Err(); err != nil {
		return model.NewError(500, "Internal Server Error", err)
	}
	return nil
}
