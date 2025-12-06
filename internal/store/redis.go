// Package store where all items regarding data storage are handled
package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(addr string) *Redis {
	if addr == "" {
		return nil
	}

	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &Redis{Client: rdb}
}

func (r *Redis) SetJTI(ctx context.Context, key, userID string, exp time.Time) error {
	return r.Client.Set(ctx, key, userID, time.Until(exp)).Err()
}

func (r *Redis) DelJTI(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *Redis) GetJTI(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}
