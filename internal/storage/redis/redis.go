package redisdb

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
)

type StorageRedis struct {
	Conn *redis.Client
	key  string
}

func New(key, port string) *StorageRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &StorageRedis{
		rdb,
		key,
	}
}

func (s *StorageRedis) GetSliceFib(ctx context.Context, x, y int) ([]string, error) {
	result, err := s.Conn.Get(ctx, s.key).Result()
	if err != nil {
		return nil, err
	}
	return strings.Split(result, ", ")[x-1 : y], nil
}

func (s *StorageRedis) SetSliceFib(ctx context.Context, fibSlice []string) error {
	err := s.Conn.Set(ctx, s.key, strings.Join(fibSlice, ", "), 0).Err()
	if err != nil {
		return err
	}
	return nil
}
