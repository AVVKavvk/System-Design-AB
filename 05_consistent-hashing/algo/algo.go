package algo

import (
	"sort"
	"sync"
)

var (
	hashRing *HashRing
	once     sync.Once
	mu       sync.RWMutex
)

type HashRing struct {
	nodes   []int            // Sorted list of virtual node hashes
	nodeMap map[int]string   // Maps virtual node hash to physical node name
	hashes  map[string][]int // Maps physical node name to virtual node hash
	userIds map[string][]string
}

func InitHashRing() *HashRing {
	once.Do(func() {
		hashRing = &HashRing{
			nodes:   make([]int, 0),
			nodeMap: make(map[int]string),
			hashes:  make(map[string][]int),
			userIds: make(map[string][]string),
		}
	})
	return hashRing
}

func GetHashRing() *HashRing {
	return hashRing
}

// AddNode adds a physical server to the ring
func (hr *HashRing) AddNode(nodeName string) {

	// Create a unique key for each virtual node
	hash := GetHashForNode(nodeName)
	mu.Lock()
	hr.nodes = append(hr.nodes, hash)
	hr.nodeMap[hash] = nodeName
	mu.Unlock()

	sort.Ints(hr.nodes)
}

// GetOwner returns the server responsible for the given key
func (hr *HashRing) GetOwner(key string) (serverName string, hashOfKey int) {
	if len(hr.nodes) == 0 {
		return "", 0
	}
	hash := GetHashForKeyForGettingOwner(key)
	// Binary search to find the first node hash >= key hash
	idx := sort.Search(len(hr.nodes), func(i int) bool {
		return hr.nodes[i] >= hash
	})
	// If we've reached the end of the ring, wrap around to the first node
	if idx == len(hr.nodes) {
		idx = 0
	}
	return hr.nodeMap[hr.nodes[idx]], hash
}

func (hr *HashRing) AddUserIdToNode(nodeName string, userId string) {

	mu.Lock()
	hr.userIds[nodeName] = append(hr.userIds[nodeName], userId)
	mu.Unlock()
}
func (hr *HashRing) DeleteNode(nodeName string) {
	// Remove from nodeMap and track which hashes to remove
	hashesToRemove := make(map[int]bool)

	hash := GetHashForNode(nodeName)
	mu.Lock()
	delete(hr.nodeMap, hash)
	mu.Unlock()
	hashesToRemove[hash] = true

	mu.Lock()
	delete(hr.hashes, nodeName)
	mu.Unlock()

	mu.Lock()
	delete(hr.userIds, nodeName)
	mu.Unlock()

	// Filter the nodes slice to remove the virtual nodes
	newNodes := make([]int, 0, len(hr.nodes))
	for _, v := range hr.nodes {
		if !hashesToRemove[v] {
			newNodes = append(newNodes, v)
		}
	}
	//OPTIONAL: But Sort the newNodes slice to ensure the ring is sorted
	sort.Ints(newNodes)

	hr.nodes = newNodes
}

func (hr *HashRing) FindTheNextNodeForNode(nodeName string) (serverName string) {
	hash := GetHashForNode(nodeName)
	idx := sort.Search(len(hr.nodes), func(i int) bool {
		// it should be > not != because it will return the same server again
		return hr.nodes[i] > hash
	})
	if idx == len(hr.nodes) {
		idx = 0
	}
	return hr.nodeMap[hr.nodes[idx]]
}

func (hr *HashRing) GetNodeInfo(nodeName string) map[string]interface{} {
	mu.RLock()
	defer mu.RUnlock()
	return map[string]interface{}{
		"nodeName":   nodeName,
		"serverHash": GetHashForNode(nodeName),
		"hashes":     hr.hashes[nodeName],
		"userIds":    hr.userIds[nodeName],
	}
}

func (hr *HashRing) GetAllNodeInfo() []map[string]interface{} {
	mu.RLock()
	defer mu.RUnlock()

	nodeInfo := make([]map[string]interface{}, 0)
	nodeMap := hr.nodeMap

	nodeInfo = append(nodeInfo, map[string]interface{}{
		"nodes": nodeMap,
	})
	return nodeInfo
}
func (hr *HashRing) AddHashToNode(nodeName string, hash int) {
	mu.Lock()
	hr.hashes[nodeName] = append(hr.hashes[nodeName], hash)
	sort.Ints(hr.hashes[nodeName])
	mu.Unlock()
}

func (hr *HashRing) GetHashesForNode(nodeName string) []int {
	mu.RLock()
	defer mu.RUnlock()
	return hr.hashes[nodeName]
}

func (hr *HashRing) GetAllTheHashLessThanOrEqualToThisHash(nodeName string, hash int) []int {

	mu.RLock()
	defer mu.RUnlock()
	// Binary search to find the first node hash >= key hash
	idx := sort.Search(len(hr.nodes), func(i int) bool {
		return hr.nodes[i] >= hash
	})
	// If we've reached the end of the ring, wrap around to the first node
	if idx == len(hr.nodes) {
		idx = 0
	}
	return hr.hashes[nodeName][idx:]
}

func (hr *HashRing) GetAllTheHashGreaterThanOrEqualToThisHash(nodeName string, hash int) []int {
	mu.RLock()
	defer mu.RUnlock()
	// Binary search to find the first node hash >= key hash
	idx := sort.Search(len(hr.nodes), func(i int) bool {
		return hr.nodes[i] >= hash
	})
	// If we've reached the end of the ring, wrap around to the first node
	if idx == len(hr.nodes) {
		idx = 0
	}
	return hr.hashes[nodeName][:idx]
}
