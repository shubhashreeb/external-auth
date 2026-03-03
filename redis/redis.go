package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type RedisConfig struct {
	PrimaryName string
	Addrs       []string
	Password    string
	DB          int
}

type RedisCache struct {
	redisConn redis.UniversalClient
}

func NewRedisCache(config *RedisConfig) (Cache, error) {
	redisOpts := &redis.UniversalOptions{
		MaxRetries: 5,
		Password:   getEnv("REDIS_PASSWORD", ""),
	}
	if config.PrimaryName != "" {
		redisOpts.MasterName = config.PrimaryName
	}
	redisOpts.Addrs = config.Addrs
	redisConn := redis.NewUniversalClient(redisOpts)
	_, err := redisConn.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to Redis. Error: %v", err)
	}
	return &RedisCache{redisConn: redisConn}, nil
}

func (c *RedisCache) Get(k string) ([]byte, error) {
	result, err := c.redisConn.Get(k).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("error getting key %s, Error : %v", k, err)
		}
		return nil, err
	}
	return []byte(result), nil
}

func (c *RedisCache) Set(k string, v []byte, d time.Duration) error {
	err := c.redisConn.Set(k, v, d).Err()
	if err != nil {
		return fmt.Errorf("error setting key %s, Error : %v", k, err)
	}
	return nil
}

func (c *RedisCache) Delete(k string) error {
	n, err := c.redisConn.Del(k).Result()
	if err != nil || n == 0 {
		return fmt.Errorf("error deleting key %s, Error: %v", k, err)
	}
	return nil
}
