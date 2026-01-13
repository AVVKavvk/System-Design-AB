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
	nodes   []int          // Sorted list of virtual node hashes
	nodeMap map[int]string // Maps virtual node hash to physical node name

}

func InitHashRing() *HashRing {
	once.Do(func() {
		hashRing = &HashRing{
			nodes:   make([]int, 0),
			nodeMap: make(map[int]string),
		}
	})
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
func (hr *HashRing) GetOwner(key string) string {
	if len(hr.nodes) == 0 {
		return ""
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
	return hr.nodeMap[hr.nodes[idx]]
}

func (hr *HashRing) DeleteNode(nodeName string) {
	// Remove from nodeMap and track which hashes to remove
	hashesToRemove := make(map[int]bool)

	hash := GetHashForNode(nodeName)
	mu.Lock()
	delete(hr.nodeMap, hash)
	mu.Unlock()
	hashesToRemove[hash] = true

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

func (hr *HashRing) FindTheNextNodeForNode(nodeName string) string {
	hash := GetHashForNode(nodeName)
	idx := sort.Search(len(hr.nodes), func(i int) bool {
		return hr.nodes[i] >= hash
	})
	if idx == len(hr.nodes) {
		idx = 0
	}
	return hr.nodeMap[hr.nodes[idx]]
}
