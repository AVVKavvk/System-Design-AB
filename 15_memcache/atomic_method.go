package main

import "log"

func atomicMethod() {

	err := mc.Delete(key1)
	if err != nil {
		log.Panicf("Failed to delete key %s", key1)
	} else {
		log.Printf("Deleted key %s", key1)
		count, err := mc.Decrement(totalUserKey, 1)
		if err != nil {
			log.Panicf("Failed to decrement key %s", totalUserKey)
		} else {

			log.Printf("Total users: %d", count)
		}

	}

}
