package main

import "log"

func retrievalMethod() {

	res, err := mc.Get(key1)
	if err != nil {
		log.Panicf("Failed to get key %s", key1)
	}
	log.Printf("User %s data: %s", key1, res.Value)

	results, err := mc.GetMulti([]string{key2, key3})
	if err != nil {
		log.Panicf("Failed to get key %s", key2)
	}
	for _, result := range results {
		log.Printf("User %s data: %s", result.Key, result.Value)
	}

	res, err = mc.GetAndTouch(key1, int32(60*10)) // Again 10 minutes
	if err != nil {
		log.Panicf("Failed to get and touch key %s", key1)
	}
	log.Printf("User %s data: %s", key1, res.Value)
}
