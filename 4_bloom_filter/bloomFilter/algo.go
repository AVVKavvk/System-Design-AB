package bloomFilter

import (
	"fmt"
	"hash/fnv"
	"sync"
)

var (
	bloomFilter *BloomFilter = nil
	once        sync.Once
)

type BloomFilter struct {
	bits          []uint64
	sizeInBits    int
	logicalColumn int
	logicalRow    int
}

func InitBloomFilter(sizeInKB int, columnSize int) *BloomFilter {
	once.Do(func() {
		totalBits := sizeInKB * 8 * 1024 // size * 8 * 1024 bits
		row := totalBits / columnSize
		bloomFilter = &BloomFilter{
			bits:          make([]uint64, row),
			sizeInBits:    totalBits,
			logicalColumn: columnSize,
			logicalRow:    row,
		}
	})
	return bloomFilter
}

func GetBloomFilter() *BloomFilter {
	return bloomFilter
}

// hash computes a single hash value
func (bf *BloomFilter) hash(data []byte) int {
	h := fnv.New64()
	_, _ = h.Write(data)
	return int(h.Sum64())
}

// Size returns the size of the bit array
func (bf *BloomFilter) Size() int {
	return bf.sizeInBits
}

// Clear resets all bits to zero
func (bf *BloomFilter) Clear() {
	for i := range bf.bits {
		bf.bits[i] = 0
	}
}

// Add inserts an element into the bloom filter
func (bf *BloomFilter) Add(data []byte) {
	index := bf.hash(data) % bf.sizeInBits
	bf.setBit(index)
}

// Contains checks if an element might be in the set
func (bf *BloomFilter) Contains(data []byte) (isFound bool, rowIdx int, colIdx int) {
	index := bf.hash(data) % bf.sizeInBits
	return bf.getBit(index)
}

func (bf *BloomFilter) setBit(pos int) (rowIdx int, colIdx int) {
	rowIndex := pos / bf.logicalColumn
	colIndex := pos % bf.logicalColumn

	fmt.Printf("Row: %d, Column: %d\n", rowIndex, colIndex)

	bf.bits[rowIndex] |= (1 << colIndex)
	return rowIndex, colIndex
}

func (bf *BloomFilter) getBit(pos int) (isFound bool, rowIdx int, colIdx int) {
	rowIndex := pos / bf.logicalColumn
	colIndex := pos % bf.logicalColumn
	return (bf.bits[rowIndex] & (1 << colIndex)) != 0, rowIndex, colIndex
}
