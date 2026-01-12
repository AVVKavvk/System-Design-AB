package helper

import (
	"hash/fnv"
	"log"
)

func GetShardIndexByClientId(userId string) int {
	hasher := fnv.New32a()
	_, err := hasher.Write([]byte(userId))
	if err != nil {
		log.Fatal("Error while hashing user id :", err)
	}

	return int(hasher.Sum32() % 3)
}
