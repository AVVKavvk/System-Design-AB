package service

import (
	"context"
	"time"

	"github.com/AVVKavvk/cache_locking/redis_client"
)

func GetAllProductsForDashboardService(c context.Context) (map[string]interface{}, error) {

	key := "redis:dashboard:product"

	ctxTimeOut, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := redis_client.GetOrUpdateProductsForDashboardWithBackoff(ctxTimeOut, key)

	if err != nil {
		return nil, err
	}

	return result, nil
}
