package algo

import "hash/crc32"

func GetHashForNode(nodeName string) int {
	hash := int(crc32.ChecksumIEEE([]byte(nodeName)))
	return hash
}
func GetHashForKeyForGettingOwner(key string) int {
	hash := int(crc32.ChecksumIEEE([]byte(key)))
	return hash
}
