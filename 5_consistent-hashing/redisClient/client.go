package redisClient

import (
	"sync"

	"github.com/go-redis/redis"
)

var (
	rc   *redis.Client = nil
	once sync.Once
)

func GetRedisClient() *redis.Client {
	return rc
}
func init() {
	once.Do(func() {
		rc = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "mypass",
			DB:       0, // use default DB
		})
	})
}
