package redis_client

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"time"

	"github.com/AVVKavvk/cache_locking/utils"
)

func GetOrUpdateProductsForDashboardWithBackoff(c context.Context, key string) (map[string]interface{}, error) {
	rdb := GetRedisClient()

	lockKey := key + ":lock"
	maxRetries := 10
	baseDelay := 50 * time.Millisecond
	maxDelay := 2 * time.Second

	log.Printf("[%s] Starting cache lookup for key: %s", getRequestID(c), key)

	for i := 0; i < maxRetries; i++ {
		log.Printf("[%s] Attempt %d/%d - Checking cache", getRequestID(c), i+1, maxRetries)

		// 1. Try to get the data with key
		dataStr := rdb.Get(key)
		if dataStr.Err() == nil {
			log.Printf("[%s] Cache HIT - Data found in cache", getRequestID(c))
			bytesData, err := dataStr.Bytes()
			if err != nil {
				log.Printf("[%s] ERROR - Failed to convert cache data to bytes: %v", getRequestID(c), err)
				return nil, err
			}
			return FromJSON[map[string]interface{}](bytesData)
		}

		log.Printf("[%s] Cache MISS - Data not found, attempting to acquire lock", getRequestID(c))

		// 2. Attempt to become the leader and set data to redis
		locked, err := rdb.SetNX(lockKey, "locked", 10*time.Second).Result()

		if err != nil {
			log.Printf("[%s] ERROR - Failed to acquire lock: %v", getRequestID(c), err)
			return nil, err
		}

		if locked {
			log.Printf("[%s] LOCK ACQUIRED - This request will update the cache", getRequestID(c))
			defer func() {
				rdb.Del(lockKey)
				log.Printf("[%s] LOCK RELEASED", getRequestID(c))
			}()

			log.Printf("[%s] Fetching data from external service...", getRequestID(c))
			time.Sleep(5 * time.Second) // Simulate backoff
			products, err := externalServiceForDashboardData()

			if err != nil {
				log.Printf("[%s] ERROR - External service failed: %v", getRequestID(c), err)
				return nil, err
			}

			log.Printf("[%s] External service returned %d products", getRequestID(c), len(products))

			bytesData, err := ToJSON(products)
			if err != nil {
				log.Printf("[%s] ERROR - Failed to marshal products to JSON: %v", getRequestID(c), err)
				return nil, err
			}

			rdb.Set(key, bytesData, 10*time.Minute)
			log.Printf("[%s] Cache UPDATED successfully (TTL: 10 minutes)", getRequestID(c))
			return products, nil
		} else {
			log.Printf("[%s] LOCK BUSY - Another request is updating the cache, waiting...", getRequestID(c))
		}

		// 3. Exponential Backoff Calculation: base * 2^retry
		delay := time.Duration(float64(baseDelay) * math.Pow(2, float64(i)))
		if delay > maxDelay {
			delay = maxDelay
		}

		log.Printf("[%s] Backing off for %v before retry", getRequestID(c), delay)

		// 4. Sleep with context awareness
		select {
		case <-c.Done():
			log.Printf("[%s] Request CANCELLED - Context done: %v", getRequestID(c), c.Err())
			return nil, c.Err()
		case <-time.After(delay):
			// Loop continues to retry
		}
	}

	log.Printf("[%s] TIMEOUT - Exceeded max retries (%d)", getRequestID(c), maxRetries)
	return nil, errors.New("request timed out waiting for cache update")
}

// Helper function to get or generate a request ID from context
func getRequestID(c context.Context) string {
	return utils.GetRequestIDFromContext(c)
}

func externalServiceForDashboardData() (map[string]interface{}, error) {
	products := map[string]interface{}{
		"p1": map[string]interface{}{
			"name":     "Product 1",
			"price":    100,
			"quantity": 10,
		},
		"p2": map[string]interface{}{
			"name":     "Product 2",
			"price":    200,
			"quantity": 20,
		},
		"p3": map[string]interface{}{
			"name":     "Product 3",
			"price":    300,
			"quantity": 30,
		},
	}

	return products, nil
}

func ToJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}

func FromJSON[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	return result, err
}
