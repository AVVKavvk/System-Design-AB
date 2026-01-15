package service

import (
	"strconv"
	"sync"

	"github.com/AVVKavvk/consistent-hashing/algo"
	"github.com/AVVKavvk/consistent-hashing/models"
	"github.com/AVVKavvk/consistent-hashing/redisClient"
)

var wg sync.WaitGroup

func AddServerService(server *models.CreateServer) map[string]interface{} {
	hr := algo.GetHashRing()
	hr.AddNode(server.Name)
	result := hr.GetNodeInfo(server.Name)
	return result
}

func GetAllServerInfoService() []map[string]interface{} {
	hr := algo.GetHashRing()
	return hr.GetAllNodeInfo()
}

func GetInfoOfServerByName(serverName string) map[string]interface{} {
	hr := algo.GetHashRing()
	return hr.GetNodeInfo(serverName)
}

func DeleteServerService(name string) map[string]interface{} {
	hr := algo.GetHashRing()
	nextServerName := hr.FindTheNextNodeForNode(name)
	allHashesForThisNode := hr.GetHashesForNode(name)

	wg.Add(len(allHashesForThisNode))

	for _, hash := range allHashesForThisNode {
		go deleteHelper(name, hash, nextServerName, &wg)
	}
	wg.Wait()
	hr.DeleteNode(name)
	return map[string]interface{}{
		"message":   "done",
		"newServer": nextServerName,
	}

}

func deleteHelper(name string, hash int, nextServerName string, wg *sync.WaitGroup) error {
	defer wg.Done()

	defer func() {
		// Recover from panic
		if r := recover(); r != nil {
			// Print stack trace
			panic(r)
		}
	}()

	hr := algo.GetHashRing()
	userData, err := redisClient.GetUserDataWithHashFromRedisWithNode(name, strconv.Itoa(hash))
	if err != nil {
		return err
	}

	err = redisClient.StoreUserDataWithHashToRedisWithNode(nextServerName, strconv.Itoa(hash), userData)

	if err != nil {
		return err
	}
	hr.AddHashToNode(nextServerName, hash)

	err = redisClient.DeleteUserDataWithHashFromRedisWithNode(name, strconv.Itoa(hash))
	if err != nil {
		return err
	}
	return nil
}
