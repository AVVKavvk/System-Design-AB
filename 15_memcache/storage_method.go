package main

import (
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

func storageMethod() {

	err := mc.Set(&memcache.Item{Key: key1, Value: []byte("User_0_data"), Expiration: int32(60 * 10)}) // 10 minutes
	if err != nil {
		log.Panicf("Failed to set key %s", key1)
	}
	res, err := mc.Increment(totalUserKey, 1)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			// Key doesn't exist, so initialize it with "1"
			// Use Add to avoid race conditions (only creates if still missing)
			err = mc.Add(&memcache.Item{Key: totalUserKey, Value: []byte("1")})
			if err != nil && err != memcache.ErrNotStored {
				log.Panicf("Initial setup failed: %v", err)
			}
			res = 1 // Our manual initial value
		} else {

			log.Panicf("Failed to increment key %s", totalUserKey)
		}
	}
	log.Printf("Total users: %d", res)

	err = mc.Set(&memcache.Item{Key: key2, Value: []byte("User_1_data")})
	if err != nil {
		log.Panicf("Failed to set key %s", key2)
	}
	res, err = mc.Increment(totalUserKey, 1)
	if err != nil {
		log.Panicf("Failed to increment key %s", totalUserKey)
	}
	log.Printf("Total users: %d", res)

	err = mc.Set(&memcache.Item{Key: key3, Value: []byte("User_2_data")})
	if err != nil {
		log.Panicf("Failed to set key %s", key3)
	}
	res, err = mc.Increment(totalUserKey, 1)
	if err != nil {
		log.Panicf("Failed to increment key %s", totalUserKey)
	}
	log.Printf("Total users: %d", res)

	// Update user_1 data

	err = mc.Replace(&memcache.Item{Key: key2, Value: []byte("User_1_new_data")})
	if err != nil {
		log.Panicf("Failed to replace key %s", key2)
	}
}
