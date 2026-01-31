package main

import (
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	addr = "localhost:11211"
	mc   *memcache.Client

	key1         = "user_0"
	key2         = "user_1"
	key3         = "user_2"
	totalUserKey = "total_users"
)

func main() {

	mc = memcache.New(addr)
	err := mc.Ping()
	if err != nil {
		log.Panicf("Failed to connect with memcached")
	} else {
		log.Printf("Connected to memcached")
	}

	storageMethod()
	retrievalMethod()
	atomicMethod()

}
