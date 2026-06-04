package redisClient

import (
	"time"

	"github.com/go-redis/redis"
)

type UrlRedis struct{}

func (ur *UrlRedis) AddUrlMapping(longUrl string, shortUrl string, ttlSeconds int) error {
	client := GetRedisClient()

	_, err := client.Set(shortUrl, longUrl, time.Duration(ttlSeconds)*time.Second).Result()
	if err != nil {
		return err
	}
	return nil
}

func (ur *UrlRedis) GetUrlMapping(shortUrl string) (string, error) {
	client := GetRedisClient()

	longUrl, err := client.Get(shortUrl).Result()
	if err == redis.Nil {
		// Key doesn't exist in cache
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return longUrl, nil
}
